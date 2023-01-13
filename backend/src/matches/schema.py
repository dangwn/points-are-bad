from pydantic import BaseModel

class Match(BaseModel):
    match_date: str
    home: str
    away: str

class MatchScore(BaseModel):
    home_goals: int
    away_goals: int

class MatchScoreWithId(MatchScore):
    match_id: int

class DisplayMatch(BaseModel):
    match_id: int
    match_date: str
    home: str
    away: str
    home_goals: int
    away_goals: int
        
    class Config:
        orm_mode = True
