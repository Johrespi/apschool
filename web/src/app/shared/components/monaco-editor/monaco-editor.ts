import {
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  OnDestroy,
  effect,
  inject,
  input,
  output,
  signal,
  viewChild,
} from '@angular/core';
import { MonacoLoaderService } from '../../../core/services/monaco-loader';

type Monaco = typeof import('monaco-editor');
type IStandaloneCodeEditor = import('monaco-editor').editor.IStandaloneCodeEditor;

@Component({
  selector: 'app-monaco-editor',
  changeDetection: ChangeDetectionStrategy.OnPush,
  template: `<div #editorContainer class="editor-container"></div>`,
  styles: `
    :host {
      display: block;
      width: 100%;
      height: 100%;
    }
    .editor-container {
      width: 100%;
      height: 100%;
    }
  `,
})
export class MonacoEditorComponent implements OnDestroy {
  private readonly monacoLoader = inject(MonacoLoaderService);

  // Inputs
  value = input<string>('');
  language = input<string>('python');
  theme = input<string>('vs-dark');
  readOnly = input<boolean>(false);

  // Outputs
  valueChange = output<string>();

  // View child
  editorContainer = viewChild.required<ElementRef<HTMLDivElement>>('editorContainer');

  // Internal state
  private editor: IStandaloneCodeEditor | null = null;
  private monaco: Monaco | null = null;
  private initialized = false;

  readonly isLoading = signal(true);

  constructor() {
    // Cargar Monaco y crear editor cuando el container esté listo
    effect(() => {
      const container = this.editorContainer();
      if (!container || this.initialized) return;

      this.initialized = true;
      this.initEditor(container.nativeElement);
    });

    // Actualizar valor cuando cambia el input (solo después de inicializado)
    effect(() => {
      const newValue = this.value();
      if (this.editor && this.editor.getValue() !== newValue) {
        this.editor.setValue(newValue);
      }
    });
  }

  private initEditor(container: HTMLElement): void {
    this.monacoLoader.load().then(monaco => {
      this.monaco = monaco;
      this.isLoading.set(false);

      this.editor = monaco.editor.create(container, {
        value: this.value(),
        language: this.language(),
        theme: this.theme(),
        readOnly: this.readOnly(),
        automaticLayout: true,
        minimap: { enabled: false },
        fontSize: 14,
        lineNumbers: 'on',
        scrollBeyondLastLine: false,
        wordWrap: 'on',
      });

      this.editor.onDidChangeModelContent(() => {
        const newValue = this.editor?.getValue() ?? '';
        this.valueChange.emit(newValue);
      });
    });
  }

  ngOnDestroy(): void {
    this.editor?.dispose();
  }
}
