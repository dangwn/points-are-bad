from fastapi import APIRouter, Depends, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.orm import Session

from authentication.schema import Token
from authentication.utils import verify_user, create_access_token
import db
from http_exceptions import USERNAME_OR_PASSWORD_EXCEPTION

router = APIRouter(
    prefix = '/login',
    tags = ['login']
)

@router.post('/', status_code = status.HTTP_202_ACCEPTED, response_model = Token)
async def login_user(
    request: OAuth2PasswordRequestForm = Depends(),
    database: Session = Depends(db.get_db)
) -> Token:
    '''
    Given a username and password, return an authentication token with a fixed lifespan
    '''
    user = await verify_user(database, request.username, request.password)
    
    if not user:
        raise USERNAME_OR_PASSWORD_EXCEPTION
    access_token = create_access_token(user_id = user.id)
    return Token(access_token = access_token, token_type = 'Bearer')
