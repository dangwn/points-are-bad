from sqlalchemy import Column, Integer, String, Date
from sqlalchemy.orm import relationship, RelationshipProperty

from db import Base
from datetime import date

class Match(Base):
    '''
    Matches table in database
    '''
    __tablename__: str = 'matches'
    
    match_id: Column = Column(Integer, primary_key=True, autoincrement=True)
    match_date: Column = Column(Date, nullable=False)
    home: Column = Column(String(64), nullable=False)
    away: Column = Column(String(64), nullable=False)
    home_goals: Column = Column(Integer, nullable=True, default=None)
    away_goals: Column = Column(Integer, nullable=True, default=None)

    # predictions: RelationshipProperty = relationship('Prediction', back_populates='match', uselist=True)

    def __init__(
        self,
        match_date: date,
        home: str,
        away: str
    ):        
        self.match_date = match_date
        self.home = home
        self.away = away
