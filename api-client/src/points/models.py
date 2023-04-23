from sqlalchemy import Column, Integer, Integer, ForeignKey
from sqlalchemy.orm import relationship, RelationshipProperty

from db import Base

class Points(Base):
    '''
    Points table in database
    '''
    __tablename__: str = 'points'

    user_id: Column = Column(Integer, ForeignKey('users.user_id', ondelete='CASCADE', onupdate='CASCADE'), primary_key=True)
    points: Column = Column(Integer, nullable=False, index=True, default=0)
    correct_scores: Column = Column(Integer, nullable=False, index=True, default=0)
    largest_error: Column = Column(Integer, nullable=False, index=True, default=0)
    position: Column = Column(Integer, nullable=True, index=True, default=None)

    user: RelationshipProperty = relationship('User', back_populates='user_points')