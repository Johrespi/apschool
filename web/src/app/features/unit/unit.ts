import { ChangeDetectionStrategy, Component, inject, input } from '@angular/core';
import { toSignal, toObservable } from '@angular/core/rxjs-interop';
import { switchMap } from 'rxjs';
import { Card } from '../../shared/components/card/card';
import { ChallengesService } from '../../core/services/challenge';

@Component({
  selector: 'app-unit',
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [Card],
  templateUrl: './unit.html',
  styleUrl: './unit.scss',
})
export class Unit {
  private readonly challengesService = inject(ChallengesService)

  slug = input.required<string>();

  challenges = toSignal(
    toObservable(this.slug).pipe(
    switchMap(slug => this.challengesService.getByCategory(slug))
    ),
    { initialValue: [] }
  )
}
