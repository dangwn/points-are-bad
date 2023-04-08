from sqlalchemy import Column, Integer, String, Boolean
from sqlalchemy.orm import relationship, RelationshipProperty

from db import Base
from config import USERNAME_MAX_LENGTH

class User(Base):
    '''
    User table in database
    '''
    __tablename__: str = 'users'

    user_id: Column = Column(Integer, primary_key=True, autoincrement=True)
    username: Column = Column(String(USERNAME_MAX_LENGTH), nullable=False)
    email: Column = Column(String(255), unique=True, nullable=False)
    hashed_password: Column = Column(String(255), nullable=False)
    is_admin: Column = Column(Boolean(), nullable = False, default=False)
    is_validated: Column = Column(Boolean(), nullable=False, default=False)

    user_points: RelationshipProperty = relationship('Points', back_populates='user', uselist=False)

    def __init__(
        self,
        username: str,
        email: str,
        hashed_password: str,
        is_admin: bool = False,
        is_validated: bool = False
    ) -> None:
        self.username: str = username
        self.email: str = email
        self.hashed_password: str = hashed_password
        self.is_admin: bool = is_admin
        self.is_validated: bool = is_validated