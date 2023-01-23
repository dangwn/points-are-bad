from sqlalchemy import Column, Integer, String, Boolean
from db import Base
import config

class User(Base):
    '''
    User model in database
    '''
    __tablename__ = 'users'

    id = Column(Integer, primary_key = True, autoincrement = True)
    display_name = Column(String(config.USERNAME_MAX_LENGTH), unique = True, nullable = False)
    email = Column(String(255), nullable = False)
    avatar = Column(String(510), nullable = False)
    provider = Column(String(255), nullable = False)
    is_admin = Column(Boolean(), nullable = False, default = False)

    def __init__(self, display_name: str, email: str, avatar: str, provider: str, is_admin: str = False) -> None:
        self.display_name = display_name
        self.email = email
        self.avatar = avatar
        self.provider = provider
        self.is_admin = is_admin
