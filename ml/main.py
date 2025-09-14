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

@app.post("/predict/")
async def predict(file: UploadFile = File(...)):
    img = Image.open(BytesIO(await file.read()))

    features = load_image(img).flatten()

    y_pred = clf.predict(features)
    label = encoder.inverse_transform(y_pred)[0]
    
    return {"prediction": label}
