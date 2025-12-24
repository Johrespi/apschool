import { inject, Injectable } from '@angular/core';
import { map, Observable } from 'rxjs';
import { Api } from './api';
import { Challenge } from '../models/challenge';

@Injectable({
  providedIn: 'root'
})
export class ChallengesService {
  private readonly api = inject(Api);

  getByCategory(category: string): Observable<Challenge[]> {
    return this.api
      .get<{challenges: Challenge[]}>(`/challenges?category=${category}`)
      .pipe(map(res => res.challenges));
  }

  getById(id : number) : Observable<Challenge> {
    return this.api
      .get<{ challenge: Challenge}>(`/challenges/${id}`)
      .pipe(map(res => res.challenge));
  }
}
