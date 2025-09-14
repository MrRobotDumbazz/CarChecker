import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-scanning',
  imports: [],
  templateUrl: './scanning.html',
  styleUrl: './scanning.css'
})
export class Scanning implements OnInit {
  constructor(private router: Router) {}

  ngOnInit(): void {
    setTimeout(() => {
      // имитация результата
      const isClean = Math.random() > 0.5;
      this.router.navigate(['/result'], { state: { clean: isClean } });
    }, 3000);
  }
}
