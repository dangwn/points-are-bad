from sqlalchemy import Column, Integer, ForeignKey
from sqlalchemy.orm import relationship, RelationshipProperty

from db import Base
from match.models import Match

class Prediction(Base):
    '''
    Predictions table in the database
    '''
    __tablename__: str = 'predictions'

    prediction_id: Column = Column(Integer, primary_key=True, autoincrement=True)
    home_goals: Column = Column(Integer, nullable=True, default=None)
    away_goals: Column = Column(Integer, nullable=True, default=None)
    user_id: Column = Column(Integer, ForeignKey('users.user_id', ondelete='CASCADE', onupdate='CASCADE'))
    match_id: Column = Column(Integer, ForeignKey(Match.match_id, ondelete='CASCADE', onupdate='CASCADE'))

    user: RelationshipProperty = relationship('User', back_populates='predictions')
    # match: RelationshipProperty = relationship('Match', back_populates='predictions', uselist=False)