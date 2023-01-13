from fastapi import APIRouter, Depends, status
from fastapi.security import OAuth2PasswordRequestForm
from sqlalchemy.orm import Session

from authentication.schema import Token
from authentication.utils import verify_user
from db import get_db
from http_exceptions import USERNAME_OR_PASSWORD_EXCEPTION

from fastapi_jwt_auth import AuthJWT

router = APIRouter(
    prefix = '/login',
    tags = ['login']
)

@router.post('/', status_code=status.HTTP_202_ACCEPTED, response_model = Token)
async def login_user(
    request: OAuth2PasswordRequestForm = Depends(),
    Authorize: AuthJWT = Depends(),
    database: Session = Depends(get_db)
) -> Token:
    user = await verify_user(request.username, request.password, database)

    if not user: 
        raise USERNAME_OR_PASSWORD_EXCEPTION

    access_token = Authorize.create_access_token(subject=user.id)
    refresh_token = Authorize.create_refresh_token(subject=user.id)

    Authorize.set_refresh_cookies(refresh_token)

    return Token(access_token=access_token, token_type='Bearer')

@router.post('/refresh', status_code=status.HTTP_202_ACCEPTED, response_model = Token)
async def refresh_access_token(
    Authorize: AuthJWT = Depends()
):
    Authorize.jwt_refresh_token_required()

    current_user = Authorize.get_jwt_subject()
    new_access_token = Authorize.create_access_token(subject=current_user)

    return Token(access_token=new_access_token, token_type='Bearer')

@router.delete('/', status_code=status.HTTP_202_ACCEPTED)
async def logout_user(
    Authorize: AuthJWT = Depends()
):
    Authorize.jwt_required()

    Authorize.unset_jwt_cookies()
    return {'msg':'Logout successful'}