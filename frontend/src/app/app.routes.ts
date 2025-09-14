import { Routes } from '@angular/router';
import { Upload } from './pages/upload/upload';
import { Scanning } from './pages/scanning/scanning';
import { Result } from './pages/result/result';
import { Final } from './pages/final/final';

export const routes: Routes = [
  { path: '', component: Upload },
  // { path: '', redirectTo: 'upload', pathMatch: 'full' },
  // { path: 'upload', component: Upload },
  { path: 'scanning', component: Scanning },
  { path: 'result', component: Result },
  { path: 'final', component: Final },
];