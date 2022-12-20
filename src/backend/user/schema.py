from pydantic import BaseModel, constr, EmailStr
import config

class User(BaseModel):
    '''
    Base class for users
    '''
    username: constr(min_length = 2, max_length = config.USERNAME_MAX_LENGTH)
    email: EmailStr
    password: str
    is_admin: bool

class DisplayUser(BaseModel):
    '''
    Class for users to be returned
    '''
    id: str
    username: str
    email: str
    is_admin: bool

    class Config:
        orm_mode = True
