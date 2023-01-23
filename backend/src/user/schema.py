from pydantic import BaseModel

class DisplayUser(BaseModel):
    display_name: str
    is_admin: bool

    class Config:
        orm_mode = True

class CreateUser(BaseModel):
    display_name: str
    avatar: str