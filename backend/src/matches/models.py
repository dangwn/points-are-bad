from sqlalchemy import Column, Integer, String
from sqlalchemy.orm import relationship
from db import Base
from datetime import datetime

class Match(Base):
    '''
    SQL Alchemy model for matches
    Match date is of form 'YYYY-MM-DD'
    '''
    __tablename__ = 'matches'
    
    match_id = Column(Integer, primary_key = True, autoincrement = True)
    match_date = Column(String(10), nullable = False) 
    home = Column(String(64), nullable = False)
    away = Column(String(64), nullable = False)
    home_goals = Column(Integer, nullable = False)
    away_goals = Column(Integer, nullable = False)

    # prediction = relationship('Prediction', back_populates = 'match')

    def __init__(
        self,
        match_date: str,
        home: str,
        away: str,
        home_goals: int = 0,
        away_goals: int = 0
    ):
        try:
            datetime.strptime(match_date, '%Y-%m-%d')
        except ValueError:
            raise ValueError('Provided date not of correct format (YYYY-MM-DD)')
        
        self.match_date = match_date
        self.home = home
        self.away = away
        self.home_goals = home_goals
        self.away_goals = away_goals