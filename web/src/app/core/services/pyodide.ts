import {Injectable, signal} from '@angular/core';

declare const loadPyodide: () => Promise<PyodideInterface>;

interface PyodideInterface {
  runPython: (code: string) => unknown;
  runPythonAsync: (code: string) => Promise<unknown>;
  loadPackagesFromImports: (code: string) => Promise<void>;
  setStdout: (options: {batched: (msg:string) => void}) => void;
  setStderr: (options: {batched: (msg:string) => void}) => void;
}

const PYODIDE_CDN = 'https://cdn.jsdelivr.net/pyodide/v0.29.0/full/pyodide.js';

export interface PythonResult {
  success: boolean;
  output: string;
  error?: string;
}

@Injectable({
  providedIn: 'root'
})
export class PyodideService {
  private pyodide: PyodideInterface | null = null;
  private loadPromise: Promise<PyodideInterface> | null = null;

  readonly isLoading = signal(false);
  readonly isReady = signal(false);

  async load(): Promise<void> {
    if (this.pyodide) return;
    if (this.loadPromise) {
      await this.loadPromise;
      return;
    }

    this.isLoading.set(true);
    this.loadPromise = new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = PYODIDE_CDN;
      script.onload = async () => {
        try {
          const pyodide = await loadPyodide();
          this.pyodide = pyodide;
          this.isReady.set(true);
          resolve(pyodide);
        } catch (err) {
          reject(err);
        } finally {
          this.isLoading.set(false);
        }
      };
      script.onerror = (err) => {
        this.isLoading.set(false);
        reject(err);
      };
      document.body.appendChild(script);
    });
    await this.loadPromise;
  }

  async runCode(userCode: string, testCode: string): Promise<PythonResult> {
    if (!this.pyodide) {
      await this.load();
    }

    const userOutput: string[] = [];
    const testOutput: string[] = [];

    try {
      // Paso 1: Ejecutar código del usuario y capturar su stdout
      this.pyodide!.setStdout({ batched: (msg) => userOutput.push(msg) });
      this.pyodide!.setStderr({ batched: (msg) => userOutput.push(`[ERROR] ${msg}`) });

      await this.pyodide!.loadPackagesFromImports(userCode);
      await this.pyodide!.runPythonAsync(userCode);

      const userOutputStr = userOutput.join('\n');

      // Paso 2: Inyectar USER_OUTPUT como variable global para el test
      this.pyodide!.runPython(`USER_OUTPUT = ${JSON.stringify(userOutputStr)}`);

      // Paso 3: Ejecutar el test (capturar su output por separado)
      this.pyodide!.setStdout({ batched: (msg) => testOutput.push(msg) });
      this.pyodide!.setStderr({ batched: (msg) => testOutput.push(`[ERROR] ${msg}`) });

      await this.pyodide!.runPythonAsync(testCode);

      const testOutputStr = testOutput.join('\n');
      const passed = testOutputStr.includes('ALL_TESTS_PASSED');

      // Combinar outputs para mostrar al usuario
      const fullOutput = userOutputStr + (testOutputStr ? '\n---\n' + testOutputStr : '');

      return {
        success: passed,
        output: fullOutput,
      };
    } catch (err) {
      const combinedOutput = [...userOutput, ...testOutput].join('\n');
      return {
        success: false,
        output: combinedOutput,
        error: err instanceof Error ? err.message : String(err),
      };
    }
  }
}
