from sqlalchemy import Column, Integer, ForeignKey
from sqlalchemy.orm import relationship
from db import Base

from matches.models import Match
from user.models import User

class Prediction(Base):
    '''
    SQL Alchemy model for predictions
    '''
    __tablename__ = 'predictions'
    
    prediction_id = Column(Integer, primary_key = True, autoincrement = True)
    user_id = Column(Integer, ForeignKey(User.id, ondelete = 'CASCADE'))
    match_id = Column(Integer, ForeignKey(Match.match_id, ondelete = 'CASCADE'))
    predicted_home_goals = Column(Integer, nullable = False)
    predicted_away_goals = Column(Integer, nullable = False)

    # match = relationship('Match', back_populates = 'prediction')
    # user = relationship('User', back_populates = 'predictions')

    def __init__(
        self,
        user_id: str,
        match_id: str,
        predicted_home_goals: int = 5,
        predicted_away_goals: int = 5
    ):
        self.user_id = user_id
        self.match_id = match_id
        self.predicted_home_goals = predicted_home_goals
        self.predicted_away_goals = predicted_away_goals