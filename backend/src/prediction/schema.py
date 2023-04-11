from pydantic import BaseModel

from match.schema import Match

class Prediction(BaseModel):
    prediction_id: int
    home_goals: int
    away_goals: int

    class Config:
        orm_mode: True

class PredictionWithMatch(Prediction):
    match: Match