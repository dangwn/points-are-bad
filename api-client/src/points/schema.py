from pydantic import BaseModel

from user.schema import Username, SessionUser

from typing import Optional

class Points(BaseModel):
    points: int
    correct_scores: int
    largest_error: int
    position: Optional[int]

    class Config:
        orm_mode: bool = True

class UserWithPoints(Points):
    user: SessionUser

class LeaderBoardUser(Points):
    user: Username