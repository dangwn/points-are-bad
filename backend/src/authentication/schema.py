from pydantic import BaseModel

class SessionUser(BaseModel):
    id: int
    display_name: str
    email: str
    avatar: str
    provider: str
    is_admin: bool

class HeaderCredentials(BaseModel):
    access_token: str
    provider: str