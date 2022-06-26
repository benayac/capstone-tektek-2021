from fastapi import Depends, FastAPI, Request
from fastapi.middleware.cors import CORSMiddleware
from requests import Session

from model import TestPredictBatch, TextComparation

import helper

from ml import TokenSimilarity


model = TokenSimilarity()
model.load_pretrained('indobenchmark/indobert-base-p2')

app = FastAPI()
# app = Flask(__name__)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/test")
def test():
    return {"message":"success"}
    
@app.post("/predict")
async def predicts(request: TextComparation):
    result = helper.predict(request.first_text, request.second_text, model)
    response = {
        'success': True,
        'request': {
            'first_text': request.first_text,
            'second_text': request.second_text
        },
        'message': 'Berhasil memprediksi kalimat',
        'accuracy': result
    }
    return result

@app.post("/predict_list")
async def predict_list(request: TestPredictBatch):
    accuracy_list = []
    for text in request.list_text:
        result = helper.predict(request.original_text, text)
        accuracy_list.append(result)
    return {
        'success': True,
        'message': 'Berhasil memprediksi kalimat',
        'accuracy': accuracy_list
    }