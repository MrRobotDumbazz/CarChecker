from fastapi import FastAPI, UploadFile, File
import joblib
import pandas as pd
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
    ImagePath: str

@app.post("/api/predict")
async def predict(request: PredictRequest):
    if not os.path.isfile(request.ImagePath):
        raise HTTPException(status_code=404, detail="Файл не найден по указанному пути")
    
    try:
        img = Image.open(request.ImagePath)
    except FileNotFoundError:
        raise HTTPException(status_code=404, detail="Изображение не найдено")
    except Exception:
        raise HTTPException(status_code=400, detail="Не удалось загрузить изображение")
    
    try:
        img_tensor = load_image(img)
        prediction_array = clf.predict(img_tensor)
        prediction = encoder.inverse_transform(prediction_array)[0]
    except Exception:
        raise HTTPException(status_code=500, detail="Ошибка при обработке изображения")
    
    return {"prediction": prediction}
