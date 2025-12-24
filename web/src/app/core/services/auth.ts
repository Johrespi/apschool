import {computed, inject, Injectable, signal} from '@angular/core';
import { Router } from '@angular/router';
import { Api } from './api';
import { environment } from '../../../environments/environment';
import { User } from '../models/user';

const TOKEN_KEY = 'token';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private readonly api = inject(Api);
  private readonly router = inject(Router);

  private readonly userSignal = signal<User | null>(null);

  readonly user = this.userSignal.asReadonly();
  readonly isAuthenticated = computed(() => this.userSignal() !== null)

  constructor() {
    const token = localStorage.getItem(TOKEN_KEY);
    if (token) {
      this.fetchCurrentUser();
    }
  }

  login(): void {
    window.location.href = `${environment.apiURL}/api/auth/github/login`;
  }

  logout(): void {
    localStorage.removeItem(TOKEN_KEY);
    this.userSignal.set(null);
    this.router.navigate(['/']);
  }

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token);
    this.fetchCurrentUser();
  }

  private fetchCurrentUser(): void {
    this.api.get<{user: User}>('/auth/me').subscribe({
      next: (response) => this.userSignal.set(response.user),
      error: () => {
        localStorage.removeItem(TOKEN_KEY);
        this.userSignal.set(null);
      },
    });
  }
}
