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
    next: (uploadRes) => {
      const imageId = uploadRes.id; // бэкенд возвращает ID загруженного файла
      this.carCheck.predictCar(imageId).subscribe({
        next: (predictRes) => {
          const predictionId = predictRes.id;
          this.router.navigate(['/scanning'], { state: { predictionId } });
        },
        error: (err) => console.error('Prediction error:', err)
      });
    },
    error: (err) => console.error('Upload error:', err)
  });
}

}
