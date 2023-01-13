from pydantic import BaseModel

class DayScore(BaseModel):
    points: int
    correct_scores: int
    largest_error: int

    class Config:
        orm_mode = True

class DayScoreWithUserId(DayScore):
    user_id: int

class PlayerPosition(DayScore):
    position: int

class PlayerPositionWithUsername(PlayerPosition):
    username: str