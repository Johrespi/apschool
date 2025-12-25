import { ChangeDetectionStrategy, Component, computed, effect, inject, input, signal } from '@angular/core';
import { toSignal, toObservable } from '@angular/core/rxjs-interop';
import { switchMap, catchError, of, tap } from 'rxjs';
import { MarkdownComponent } from 'ngx-markdown';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBar } from '@angular/material/snack-bar';
import { ChallengesService } from '../../core/services/challenge';
import { PyodideService, PythonResult } from '../../core/services/pyodide';
import { SubmissionService } from '../../core/services/submission';
import { MonacoEditorComponent } from '../../shared/components/monaco-editor/monaco-editor';

@Component({
  selector: 'app-challenge',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [MonacoEditorComponent, MarkdownComponent, MatButtonModule],
  templateUrl: './challenge.html',
  styleUrl: './challenge.scss',
})
export class Challenge {
  private readonly challengesService = inject(ChallengesService);
  private readonly pyodideService = inject(PyodideService);
  private readonly submissionService = inject(SubmissionService);
  private readonly snackBar = inject(MatSnackBar);

  // Route params
  slug = input.required<string>();
  id = input.required<string>();

  // Fetch challenge
  challenge = toSignal(
    toObservable(this.id).pipe(
      switchMap(id => this.challengesService.getById(Number(id)))
    )
  );

  // Fetch existing submission (if any)
  existingSubmission = toSignal(
    toObservable(this.id).pipe(
      switchMap(id =>
        this.submissionService.getByChallenge(Number(id)).pipe(
          catchError(() => of(null))
        )
      ),
      tap(submission => {
        if (submission) {
          this.isSubmitted.set(true);
        }
      })
    )
  );

  // State
  code = signal('');
  output = signal('');
  isRunning = signal(false);
  isSubmitting = signal(false);
  isSubmitted = signal(false);
  lastResult = signal<PythonResult | null>(null);

  // Computed
  pyodideLoading = this.pyodideService.isLoading;
  pyodideReady = this.pyodideService.isReady;
  canSubmit = computed(() => this.lastResult()?.success === true);

  // Initialize code with existing submission or template
  initialCode = computed(() => {
    const submission = this.existingSubmission();
    if (submission?.code) return submission.code;
    return this.challenge()?.template ?? '';
  });

  submitButtonText = computed(() => {
    if (this.isSubmitting()) return 'Enviando...';
    if (this.isSubmitted()) return 'Enviado';
    return 'Enviar';
  });

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
    const challenge = this.challenge();
    if (!challenge) return;

    this.isSubmitting.set(true);

    this.submissionService
      .create({
        challenge_id: challenge.id,
        code: this.code(),
        passed: true,
      })
      .subscribe({
        next: () => {
          this.isSubmitting.set(false);
          this.isSubmitted.set(true);
          this.snackBar.open('¡Solución guardada correctamente!', 'Cerrar', {
            duration: 5000,
          });
        },
        error: () => {
          this.isSubmitting.set(false);
          this.snackBar.open('Error al guardar la solución', 'Cerrar', {
            duration: 5000,
          });
        },
      });
  }

  onCodeChange(newCode: string): void {
    this.code.set(newCode);
    // Reset result when code changes
    this.lastResult.set(null);
    // Allow resubmit if code changes after submission
    if (this.isSubmitted()) {
      this.isSubmitted.set(false);
    }
  }
}
