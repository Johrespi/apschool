import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root',
})
export class Api {
  private readonly http = inject(HttpClient);
  private readonly baseURL = 'http://localhost:8080/api';

  get<T>(path: string): Observable<T> {
    return this.http.get<T>(`${this.baseURL}${path}`);
  }

  post<T>(path:string, body: unknown): Observable<T> {
    return this.http.post<T>(`${this.baseURL}${path}`, body)
  }

  put<T>(path: string, body: unknown): Observable<T> {
    return this.http.put<T>(`${this.baseURL}${path}`, body);
  }

  patch<T>(path: string, body: unknown): Observable<T> {
    return this.http.patch<T>(`${this.baseURL}${path}`, body);
  }

  delete<T>(path: string): Observable<T> {
    return this.http.delete<T>(`${this.baseURL}${path}`);
  }
}
