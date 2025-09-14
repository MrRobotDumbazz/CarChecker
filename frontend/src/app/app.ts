import { Component, signal } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { CarCheckService } from './services/car-check.service';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterOutlet, CommonModule],
  templateUrl: './app.html',
  styleUrl: './app.css'
})
export class App {
  protected readonly title = signal('car-condition-checker');

  selectedFile: File | null = null;
  result: any;

  constructor(private carCheck: CarCheckService) {}

  onFileSelected(event: any) {
    this.selectedFile = event.target.files[0];
  }

  onUpload() {
    if (this.selectedFile) {
      this.carCheck.uploadCarPhoto(this.selectedFile).subscribe({
        next: (res) => (this.result = res),
        error: (err) => console.error('Upload error:', err),
      });
    }
  }
}
