import { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: 'auth/callback',
    loadComponent: () =>
      import('./features/auth/callback').then(m => m.AuthCallback),
  },
  {
    path: '',
    loadComponent: () =>
      import('./shared/layouts/main-layout/main-layout').then(m => m.MainLayout),
    children: [
      {
        path: '',
        loadComponent: () => import('./features/home/home').then(m => m.Home),
      },
      {
        path: 'units/:slug',
        loadComponent: () => import ('./features/unit/unit').then(m => m.Unit)
      }
    ],
  },
];
