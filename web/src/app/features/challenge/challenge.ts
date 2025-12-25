import { ChangeDetectionStrategy, Component, computed, inject, input, signal } from '@angular/core';
import { toSignal, toObservable } from '@angular/core/rxjs-interop';
import { switchMap } from 'rxjs';
import { MarkdownComponent } from 'ngx-markdown';
import { ChallengesService } from '../../core/services/challenge';
import { PyodideService, PythonResult } from '../../core/services/pyodide';
import { MonacoEditorComponent } from '../../shared/components/monaco-editor/monaco-editor';
@Component({
  selector: 'app-challenge',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MonacoEditorComponent, MarkdownComponent],
  templateUrl: './challenge.html',
  styleUrl: './challenge.scss',
})
export class Challenge {
  private readonly challengesService = inject(ChallengesService);
  private readonly pyodideService = inject(PyodideService);
  // Route params
  slug = input.required<string>();
  id = input.required<string>();
  // Fetch challenge
  challenge = toSignal(
    toObservable(this.id).pipe(
      switchMap(id => this.challengesService.getById(Number(id)))
    )
  );
  // State
  code = signal('');
  output = signal('');
  isRunning = signal(false);
  lastResult = signal<PythonResult | null>(null);
  // Computed
  pyodideLoading = this.pyodideService.isLoading;
  pyodideReady = this.pyodideService.isReady;
  canSubmit = computed(() => this.lastResult()?.success === true);
  // Initialize code with template
  initialCode = computed(() => this.challenge()?.template ?? '');
  async onRun(): Promise<void> {
    const challenge = this.challenge();
    if (!challenge) return;
    this.isRunning.set(true);
    this.output.set('Ejecutando...');
    try {
      const result = await this.pyodideService.runCode(this.code(), challenge.test_code);
      this.lastResult.set(result);

      if (result.success) {
        this.output.set(`${result.output}\n\n✅ ¡Todos los tests pasaron!`);
      } else if (result.error) {
        this.output.set(`${result.output}\n\n❌ Error: ${result.error}`);
      } else {
        this.output.set(`${result.output}\n\n❌ Los tests no pasaron`);
      }
    } catch (err) {
      this.output.set(`Error inesperado: ${err}`);
    } finally {
      this.isRunning.set(false);
    }
  }
  onSubmit(): void {
    // TODO: Llamar a POST /api/submissions
    console.log('Submit:', this.code());
  }
  onCodeChange(newCode: string): void {
    this.code.set(newCode);
    // Reset result when code changes
    this.lastResult.set(null);
  }
}
