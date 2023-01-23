from sqlalchemy import Column, Integer, ForeignKey
from sqlalchemy.orm import relationship
from db import Base
from user.models import User

class PlayerPoints(Base):
    '''
    SQL Alchemy model for player points
    '''
    __tablename__ = 'points'
    
    user_id = Column(Integer, ForeignKey(User.id, ondelete = 'CASCADE'), primary_key = True)
    points = Column(Integer, nullable = False)
    correct_scores = Column(Integer, nullable = False)
    largest_error = Column(Integer, nullable = False)

    # user = relationship('User', back_populates = 'score')

    def __init__(
        self,
        user_id: int,
        points: int = 0,
        correct_scores: int = 0,
        largest_error: int = 0
    ):
        self.user_id = user_id
        self.points = points
        self.correct_scores = correct_scores
        self.largest_error = largest_error