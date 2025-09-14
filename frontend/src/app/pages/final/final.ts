import { Component, OnInit } from '@angular/core';
import { AfterViewInit } from '@angular/core';
import * as L from 'leaflet';
import { SafeUrlPipe } from './safe-url.pipe';
import { CarCheckService } from '../../services/car-check.service';

@Component({
  selector: 'app-final',
  imports: [SafeUrlPipe],
  templateUrl: './final.html',
  styleUrl: './final.css'
})
export class Final implements OnInit {
  // ngAfterViewInit() {
  //   if (navigator.geolocation) {
  //     navigator.geolocation.getCurrentPosition((pos) => {
  //       const coords = [pos.coords.latitude, pos.coords.longitude] as [number, number];
  //       const map = L.map('map').setView(coords, 14);

  //       L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  //         attribution: '© OpenStreetMap contributors'
  //       }).addTo(map);

  //       L.marker(coords).addTo(map).bindPopup('You are here').openPopup();
  //     });
  //   }
  // }

  mapUrl: string = '';
  predictions: any[] = [];
  loading = false;

  constructor(private carCheckService: CarCheckService) {}

  ngOnInit() {
    // Получаем predictions из router state
    if (history.state.predictions) {
      this.predictions = history.state.predictions;
    }

    // Локация
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition((pos) => {
        const lat = pos.coords.latitude;
        const lng = pos.coords.longitude;
        this.mapUrl = `https://www.google.com/maps?q=${lat},${lng}&z=15&output=embed`;
      });
    }
  }

}
