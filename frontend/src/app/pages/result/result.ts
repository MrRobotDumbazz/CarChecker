import { Component } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-result',
  imports: [],
  templateUrl: './result.html',
  styleUrl: './result.css'
})
export class Result {
  isClean = history.state.clean;
  // isClean = true;

  constructor(private router: Router) {}

  goFinal() {
    this.router.navigate(['/final']);
  }
}
