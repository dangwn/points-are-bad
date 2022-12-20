from pydantic import BaseModel

class LoginRequest(BaseModel):
    '''
    Schema for the login request
    '''
    username: str
    password: str
    
class Token(BaseModel):
    '''
    Schema for the authentication token return after login
    '''
    access_token: str
    token_type: str