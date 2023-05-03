from pydantic import BaseModel, constr, EmailStr
from config import USERNAME_MIN_LENGTH, USERNAME_MAX_LENGTH

class User(BaseModel):
    user_id: int
    username: constr(min_length=USERNAME_MIN_LENGTH, max_length=USERNAME_MAX_LENGTH)
    email: EmailStr
    hashed_password: str
    is_admin: bool

    class Config:
        orm_mode: bool = True

class NewUser(BaseModel):
    token: str
    username: str
    password: str

class Username(BaseModel):
    username: constr(min_length=USERNAME_MIN_LENGTH, max_length=USERNAME_MAX_LENGTH)

    class Config:
        orm_mode: bool = True

class SessionUser(Username):
    is_admin: bool
