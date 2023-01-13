from sqlalchemy import Column, Integer, String, Boolean
from sqlalchemy.orm import relationship
from db import Base
import config

class User(Base):
    '''
    User model in database
    '''
    __tablename__ = 'users'

    id = Column(Integer, primary_key = True, autoincrement = True)
    username = Column(String(config.USERNAME_MAX_LENGTH), unique = True, nullable = False)
    email = Column(String(255), unique = True, nullable = False)
    password = Column(String(255), nullable = False)
    is_admin = Column(Boolean(), nullable = False, default = False)

    score = relationship('PlayerPoints', back_populates = 'user')
    # predictions = relationship('Prediction', back_populates = 'user')

    def __init__(self, username: str, email: str, password: str, is_admin: str = False) -> None:
        self.username = username
        self.email = email
        self.password = password
        self.is_admin = is_admin