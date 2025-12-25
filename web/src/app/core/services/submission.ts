import { inject, Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Api } from './api';
import { Submission } from '../models/submission';

@Injectable({
  providedIn: 'root',
})
export class SubmissionService {
  private readonly api = inject(Api);

  create(submission: Submission): Observable<Submission> {
    return this.api
      .post<{ submission: Submission }>('/submissions', submission)
      .pipe(map((res) => res.submission));
  }

  getByChallenge(challengeId: number): Observable<Submission | null> {
    return this.api
      .get<{ submission: Submission }>(`/submissions/${challengeId}`)
      .pipe(map((res) => res.submission));
  }
}
