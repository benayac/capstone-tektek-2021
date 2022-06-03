from typing import List, Optional
from pydantic import BaseModel

class TextComparation(BaseModel):
    first_text: str
    second_text: str

class ReportRequest(BaseModel):
    jss_id: int
    jenis: str
    keterangan: str
    latitude: float
    longitude: float
    tanggal_laporan: str

class TestPredictBatch(BaseModel):
    original_text: str
    list_text: list[str]