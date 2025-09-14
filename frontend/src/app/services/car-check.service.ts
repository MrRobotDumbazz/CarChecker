import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";

@Injectable({
  providedIn: 'root'
})
export class CarCheckService {
  // Старый:
  // private apiUrl = 'http://localhost:8080/api/check';
  
  // Новый — совпадает с бэкенд эндпоинтами:
  private uploadUrl = 'http://localhost:8080/api/v1/images/upload';

  constructor(private http: HttpClient) {}

  uploadCarPhoto(file: File): Observable<any> {
    const formData = new FormData();
    formData.append('file', file);

    return this.http.post<any>(this.uploadUrl, formData);
  }

  // Добавим метод для запуска предсказания по image_id
  predictCar(imageId: string): Observable<any> {
    return this.http.post<any>(`http://localhost:8080/api/v1/predict/${imageId}`, {});
  }

  // Получение результатов анализа
  getPrediction(predictionId: string): Observable<any> {
    return this.http.get<any>(`http://localhost:8080/api/v1/predictions/${predictionId}`);
  }
}
