from pydantic import BaseModel

class MatchPrediction(BaseModel):
    match_id: int
    predicted_home_goals: int
    predicted_away_goals: int

    class Config:
        orm_mode = True

class MatchPredictionWithUserId(MatchPrediction):
    user_id: int

class MatchPredictionWithTeams(MatchPrediction):
    home: str
    away: str