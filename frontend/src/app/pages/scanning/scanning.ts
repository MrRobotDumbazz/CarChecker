import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { CarCheckService } from '../../services/car-check.service';

@Component({
  selector: 'app-scanning',
  imports: [],
  templateUrl: './scanning.html',
  styleUrl: './scanning.css'
})
export class Scanning implements OnInit {

  constructor(private router: Router, private carCheck: CarCheckService) {}

  ngOnInit(): void {
    const predictionId = history.state.predictionId;
    if (!predictionId) return;

    const interval = setInterval(() => {
      this.carCheck.getPrediction(predictionId).subscribe({
        next: (res: { status: string }) => {
          if (res.status === 'completed' || res.status === 'failed') {
            clearInterval(interval);
            this.router.navigate(['/result'], { state: { result: res } });
          }
        },
        error: (err: any) => {
          console.error(err);
          clearInterval(interval);
        }
      });
    }, 2000);
  }
}
