from pydantic import BaseModel

from match.schema import Match
from user.schema import SessionUser
from typing import Optional

class Prediction(BaseModel):
    prediction_id: int
    home_goals: Optional[int]
    away_goals: Optional[int]

    class Config:
        orm_mode: bool = True

class PredictionWithMatch(Prediction):
    match: Match
    user: SessionUser