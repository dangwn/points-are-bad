from sqlalchemy import Column, Integer, ForeignKey, String
from db import Base
from datetime import datetime

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
    match_date = Column(String(10), nullable = False)

    def __init__(
        self,
        user_id: str,
        match_id: str,
        match_date: str, 
        predicted_home_goals: int = 5,
        predicted_away_goals: int = 5
    ):
        
        try:
            datetime.strptime(match_date, '%Y-%m-%d')
        except ValueError:
            raise ValueError('Provided date not of correct format (YYYY-MM-DD)')

        self.user_id = user_id
        self.match_id = match_id
        self.match_date = match_date
        self.predicted_home_goals = predicted_home_goals
        self.predicted_away_goals = predicted_away_goals