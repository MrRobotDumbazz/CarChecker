import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { CarCheckService } from '../../services/car-check.service';

@Component({
  selector: 'app-upload',
  standalone: true,
  imports: [],
  templateUrl: './upload.html',
  styleUrl: './upload.css'
})
export class Upload {
  // constructor(private router: Router) {}
  constructor(private router: Router, private carCheck: CarCheckService) {}
  onFileSelected(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;

  this.carCheck.uploadCarPhoto(file).subscribe({
    next: (res) => {
      console.log('Backend response:', res);
      // передаем predictions через router state
      this.router.navigate(['/final'], { state: { predictions: res.predictions } });
    },
    error: (err) => {
      console.error('Upload error:', err);
    },
  });
}

}
