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

    const output: string[] = [];

    this.pyodide!.setStdout({ batched: (msg) => output.push(msg) });
    this.pyodide!.setStderr({ batched: (msg) => output.push(`[ERROR] ${msg}`)});

    try {
      await this.pyodide!.loadPackagesFromImports(userCode);

      const fullCode = `${userCode}\n\n${testCode}`;
      await this.pyodide!.runPythonAsync(fullCode);

      const outputStr = output.join('\n');
      const passed = outputStr.includes('ALL_TESTS_PASSED');

      return {
        success: passed,
        output: outputStr,
      };
    } catch (err) {
      return {
        success: false,
        output: output.join('\n'),
        error: err instanceof Error ? err.message : String(err),
      };
    }
  }
}
