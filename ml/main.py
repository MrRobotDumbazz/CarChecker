from fastapi import FastAPI, UploadFile, File, HTTPException
from pydantic import BaseModel
import joblib
import pandas as pd
import os
from io import BytesIO
from PIL import Image
import numpy as np

app = FastAPI()

clf = joblib.load("model.pkl")
encoder = joblib.load("encoder.pkl")

def load_image(image):
    img = Image.open(image).convert("RGB").resize((640, 640))
    img = np.array(img).flatten()
    return img
    
class PredictRequest(BaseModel):
    image_path: str

@app.get("/health")
async def health():
    return {"status": "ok", "service": "ml-service"}

@app.post("/api/predict")
async def predict(request: PredictRequest):
    if not os.path.isfile(request.image_path):
        raise HTTPException(status_code=404, detail="Файл не найден по указанному пути")

    try:
        img = Image.open(request.image_path)
    except FileNotFoundError:
        raise HTTPException(status_code=404, detail="Изображение не найдено")
    except Exception:
        raise HTTPException(status_code=400, detail="Не удалось загрузить изображение")

    try:
        img_tensor = load_image(request.image_path)
        prediction_array = clf.predict([img_tensor])
        prediction = encoder.inverse_transform(prediction_array)[0]
        print(f"DEBUG: Raw prediction = {prediction}")  # Debug output
    except Exception as e:
        print(f"DEBUG: Prediction error = {e}")
        raise HTTPException(status_code=500, detail="Ошибка при обработке изображения")
    
    # Map prediction to expected backend format
    if prediction == "scratch":
        cleanliness_status = "dirty"
        integrity_status = "damaged"
    elif prediction == "car":
        cleanliness_status = "clean"
        integrity_status = "intact"
    else:
        cleanliness_status = "unknown"
        integrity_status = "unknown"

    return {
        "cleanliness": {
            "status": cleanliness_status,
            "confidence": 0.85  # Default confidence
        },
        "integrity": {
            "status": integrity_status,
            "confidence": 0.85  # Default confidence
        },
        "processing_time_ms": 100,
        "model_version": "v1.0",
        "success": True
    }
