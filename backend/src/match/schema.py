from pydantic import BaseModel

from typing import Optional
from datetime import date

class MatchWithoutGoals(BaseModel):
    match_date: date
    home: str
    away: str

    class Config:
        orm_mode: bool = True

class Match(MatchWithoutGoals):
    home_goals: Optional[int]
    away_goals: Optional[int]
