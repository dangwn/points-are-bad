from fastapi import (
    APIRouter, 
    Response, 
    Request,
    HTTPException, 
    status,
    Depends
)
from datetime import timedelta
from fastapi_csrf_protect import CsrfProtect
from sqlalchemy.orm import Session
from jose import jwt
import json

from auth.schema import LoginUser, Token
from auth.validate import validate_login_user
from auth.csrf import CsrfSettings
from auth.utils import (
    generate_jwt_token,
    create_verification_token,
    is_user_email_in_db,
    validate_email
)
from config import (
    ACCESS_TOKEN_LIFETIME_MINUTES,
    ACCESS_TOKEN_SECRET,
    REFRESH_TOKEN_LIFETIME_DAYS,
    REFRESH_TOKEN_SECRET,
    REFRESH_TOKEN_COOKIE_KEY,
    CSRF_TOKEN_COOKIE_KEY,
    FRONTEND_URL
)
from user.schema import User
from db import get_db

from typing import (
    Optional,
    Dict
)

router: APIRouter = APIRouter(
    prefix='/auth'
)

@CsrfProtect.load_config
def get_csrf_confg() -> CsrfSettings:
    '''
    Loads the CSRF settings into module
    '''
    return CsrfSettings()

@router.post('/login/', status_code=status.HTTP_202_ACCEPTED)
async def login(
    login_user: LoginUser,
    response: Response,
    csrf_protect: CsrfProtect = Depends(),
    db: Session = Depends(get_db)
) -> Token:
    '''
    Post request endpoint to log user in
    Validates username and password, and provides an access token as a response
        and a csrf and refresh token as cookies
    See config for token expiry
    '''
    user: Optional[User] = await validate_login_user(
        login_user=login_user,
        db=db
    )
    if user is None:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='Could not find user with email and password combination'
        )

    # Generate tokens
    access_token: str = generate_jwt_token(
        subject=str(user.user_id),
        expire_time=timedelta(minutes=ACCESS_TOKEN_LIFETIME_MINUTES),
        secret_key=ACCESS_TOKEN_SECRET
    )
    refresh_token: str = generate_jwt_token(
        subject=str(user.user_id),
        expire_time=timedelta(days=REFRESH_TOKEN_LIFETIME_DAYS),
        secret_key=REFRESH_TOKEN_SECRET
    )
    csrf_token: str = csrf_protect.generate_csrf()

    # Set cookies
    response.set_cookie(
        key=REFRESH_TOKEN_COOKIE_KEY,
        value=refresh_token,
        secure=True
    )
    response.set_cookie(
        key=CSRF_TOKEN_COOKIE_KEY,
        value=csrf_token,
        secure=True
    )
    return Token(access_token=access_token, token_type='Bearer')

@router.delete('/login/',status_code=status.HTTP_204_NO_CONTENT)
async def logout(
    response: Response
) -> Response:
    '''
    Delete endpoint to log user out by deleting refresh and CSRF
        token cookies 
    '''
    response.delete_cookie(REFRESH_TOKEN_COOKIE_KEY)
    response.delete_cookie(CSRF_TOKEN_COOKIE_KEY)
    response.status_code = status.HTTP_204_NO_CONTENT
    return response

@router.post('/refresh/')
async def refresh_access_token(
    request: Request,
    csrf_protect: CsrfProtect = Depends()
) -> Token:
    '''
    Post request endpoint for refreshing user access token
    Validates user's csrf token, then uses refresh token to generate a new
        access token
    '''
    refresh_token: str = request.cookies.get('X-Refresh-Token')
    csrf_token: str = request.cookies.get('X-CSRF-Token')
    
    try:
        csrf_protect.validate_csrf(csrf_token)
    except:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='Could not validate csrf token'
        )

    user_id: str = jwt.decode(
        token=refresh_token,
        key=REFRESH_TOKEN_SECRET
    )['sub']

    access_token: str = generate_jwt_token(
        subject=user_id,
        secret_key=ACCESS_TOKEN_SECRET,
        expire_time=timedelta(minutes=ACCESS_TOKEN_LIFETIME_MINUTES)
    )

    return Token(access_token=access_token, token_type='Bearer')
    
@router.post('/verify/', status_code=status.HTTP_202_ACCEPTED)
async def create_email_verification_code(
    email: str,
    request: Request,
    db: Session = Depends(get_db)
) -> Response:
    email_is_valid: bool = await validate_email(email=email)
    if not email_is_valid:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Email address format not valid'
        )
    
    email_in_db: bool = await is_user_email_in_db(
        email=email,
        db=db
    )
    if email_in_db:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail='Email address already in use'
        )
    
    verification_token: str = await create_verification_token(email=email)

    # Send email
    await request.app.rabbitmq_producer.send_message(
        json.dumps({
            email: f'{FRONTEND_URL}/signup?token={verification_token}'
        }), 
        queue_name='email_client_queue'
    )
