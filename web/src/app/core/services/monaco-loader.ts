import { Injectable, signal } from '@angular/core';

type Monaco = typeof import('monaco-editor');

const MONACO_VERSION = '0.53.0';
const MONACO_BASE_URL = `https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/${MONACO_VERSION}/min`;

@Injectable({ providedIn: 'root' })
export class MonacoLoaderService {
  private loadPromise: Promise<Monaco> | null = null;

  readonly isLoaded = signal(false);

  load(): Promise<Monaco> {
    if (this.loadPromise) {
      return this.loadPromise;
    }

    this.loadPromise = new Promise((resolve, reject) => {
      // Si ya está cargado, resolver inmediatamente
      const win = window as unknown as { monaco?: Monaco };
      if (win.monaco) {
        this.isLoaded.set(true);
        resolve(win.monaco);
        return;
      }

      const script = document.createElement('script');
      script.src = `${MONACO_BASE_URL}/vs/loader.min.js`;
      script.onload = () => {
        const require = (window as unknown as {
          require: {
            config: (options: { paths: Record<string, string> }) => void;
            (deps: string[], cb: () => void): void;
          }
        }).require;

        require.config({ paths: { vs: `${MONACO_BASE_URL}/vs` } });
        require(['vs/editor/editor.main'], () => {
          this.isLoaded.set(true);
          resolve((window as unknown as { monaco: Monaco }).monaco);
        });
      };
      script.onerror = reject;
      document.body.appendChild(script);
    });

    return this.loadPromise;
  }
}
