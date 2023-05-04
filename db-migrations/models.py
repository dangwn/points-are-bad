from sqlalchemy import (
    Boolean,
    Column, 
    Date, 
    ForeignKey,
    Integer,
    String
)
from sqlalchemy.orm import relationship, RelationshipProperty

from db import Base
from config import USERNAME_MAX_LENGTH

__all__ = [
    'Match', 'Points', 'Prediction', 'User'
]

class Match(Base):
    '''
    Matches table in database
    '''
    __tablename__: str = 'matches'
    
    match_id: Column = Column(Integer, primary_key=True, autoincrement=True)
    match_date: Column = Column(Date, nullable=False, index=True)
    home: Column = Column(String(64), nullable=False)
    away: Column = Column(String(64), nullable=False)
    home_goals: Column = Column(Integer, nullable=True, default=None)
    away_goals: Column = Column(Integer, nullable=True, default=None)

    predictions: RelationshipProperty = relationship('Prediction', back_populates='match', uselist=True)

class Points(Base):
    '''
    Points table in database
    '''
    __tablename__: str = 'points'

    user_id: Column = Column(String(36), ForeignKey('users.user_id', ondelete='CASCADE', onupdate='CASCADE'), primary_key=True)
    points: Column = Column(Integer, nullable=False, index=True, default=0)
    correct_scores: Column = Column(Integer, nullable=False, index=True, default=0)
    largest_error: Column = Column(Integer, nullable=False, index=True, default=0)
    position: Column = Column(Integer, nullable=True, index=True, default=None)

    user: RelationshipProperty = relationship('User', back_populates='user_points')

class Prediction(Base):
    '''
    Predictions table in the database
    '''
    __tablename__: str = 'predictions'

    prediction_id: Column = Column(Integer, primary_key=True, autoincrement=True)
    home_goals: Column = Column(Integer, nullable=True, default=None)
    away_goals: Column = Column(Integer, nullable=True, default=None)
    user_id: Column = Column(String(36), ForeignKey('users.user_id', ondelete='CASCADE', onupdate='CASCADE'))
    match_id: Column = Column(Integer, ForeignKey(Match.match_id, ondelete='CASCADE', onupdate='CASCADE'))
    
    user: RelationshipProperty = relationship('User', back_populates='predictions')
    match: RelationshipProperty = relationship('Match', back_populates='predictions', uselist=False)

class User(Base):
    '''
    User table in database
    '''
    __tablename__: str = 'users'

    user_id: Column = Column(String(36), primary_key=True)
    username: Column = Column(String(USERNAME_MAX_LENGTH), nullable=False)
    email: Column = Column(String(255), unique=True, nullable=False)
    hashed_password: Column = Column(String(255), nullable=False)
    is_admin: Column = Column(Boolean(), nullable = False, default=False)

    user_points: RelationshipProperty = relationship('Points', back_populates='user', uselist=False)
    predictions: RelationshipProperty = relationship('Prediction', back_populates='user', uselist=True)