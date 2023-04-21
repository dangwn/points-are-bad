from fastapi import (
    APIRouter,
    Depends,
    status,
    Response,
    HTTPException,
    Body
)
from fastapi_csrf_protect import CsrfProtect
from datetime import timedelta
from sqlalchemy.orm import Session

from db import get_db
from auth.utils import (
    hash_password,
    generate_jwt_token,
    get_current_user,
    validate_verification_token
)
from auth.schema import Token
from auth.csrf import CsrfSettings
from config import (
    ACCESS_TOKEN_LIFETIME_MINUTES,
    ACCESS_TOKEN_SECRET,
    REFRESH_TOKEN_LIFETIME_DAYS,
    REFRESH_TOKEN_SECRET,
    REFRESH_TOKEN_COOKIE_KEY,
    CSRF_TOKEN_COOKIE_KEY
)
from exceptions import (
    INCORRECT_PASSWORD_EXCEPTION,
    NO_CURRENT_USER_EXCEPTION,
    COULD_NOT_UPDATE_EXCEPTION
)
from user.models import User as UserModel
from user.schema import NewUser, SessionUser
from user.utils import (
    insert_user_into_db,
    delete_user_by_id,
    change_username_by_id,
    change_password_by_id,
    verify_current_password
)

from typing import (
    Optional
)

router: APIRouter = APIRouter(
    prefix='/user'
)

@CsrfProtect.load_config
def get_csrf_confg() -> CsrfSettings:
    '''
    Loads the CSRF settings into module
    '''
    return CsrfSettings()

@router.get('/', response_model=SessionUser)
async def display_current_user(
    current_user: UserModel = Depends(get_current_user)
) -> SessionUser:
    if not current_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='No current user'
        )
    
    return current_user


@router.post('/', status_code=status.HTTP_201_CREATED)
async def create_new_user(
    new_user: NewUser,
    response: Response,
    db: Session = Depends(get_db),
    csrf_protect: CsrfProtect = Depends()
) -> Token:
    '''
    Post request endpoint for creating a new, non-verified user
        in the db
    '''
    email: Optional[str] = await validate_verification_token(
        token=new_user.token
    )
    if not email:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='Invalid verification token'
        )

    hashed_password: str = hash_password(new_user.password)

    new_user: UserModel = await insert_user_into_db(
        username=new_user.username,
        email=email,
        hashed_password=hashed_password,
        db=db
    )
    if not new_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not create new user'
        )
    
    # Generate tokens
    access_token: str = generate_jwt_token(
        subject=str(new_user.user_id),
        expire_time=timedelta(minutes=ACCESS_TOKEN_LIFETIME_MINUTES),
        secret_key=ACCESS_TOKEN_SECRET
    )
    refresh_token: str = generate_jwt_token(
        subject=str(new_user.user_id),
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

@router.delete('/', status_code=status.HTTP_204_NO_CONTENT)
async def delete_current_user(
    response: Response,
    current_user: UserModel = Depends(get_current_user),
    db: Session = Depends(get_db)
) -> Response:
    if not current_user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='No current user'
        )
    
    await delete_user_by_id(
        user_id=current_user.user_id,
        db=db
    )

    response.delete_cookie(REFRESH_TOKEN_COOKIE_KEY)
    response.delete_cookie(CSRF_TOKEN_COOKIE_KEY)
    response.status_code = status.HTTP_204_NO_CONTENT
    return response

@router.put('/username')
async def edit_username(
    username: str = Body(),
    current_user: UserModel = Depends(get_current_user),
    db: Session = Depends(get_db)
) -> Response:
    if not current_user:
        raise NO_CURRENT_USER_EXCEPTION
    
    new_username: bool = await change_username_by_id(
        user_id=current_user.user_id,
        new_username=username,
        db=db
    )
    if not new_username:
        raise COULD_NOT_UPDATE_EXCEPTION(what="user's username")


@router.put('/password')
async def edit_password(
    new_password: str = Body(),
    current_password: str = Body(),
    current_user: UserModel = Depends(get_current_user),
    db: Session = Depends(get_db)
) -> Response:
    if not current_user:
        raise NO_CURRENT_USER_EXCEPTION
    
    current_password_match: bool = await verify_current_password(
        user_id=current_user.user_id,
        current_password=current_password,
        db=db
    )
    if not current_password_match:
        raise INCORRECT_PASSWORD_EXCEPTION
    
    new_password_created: bool = await change_password_by_id(
        user_id=current_user.user_id,
        password=new_password,
        db=db
    )
    if not new_password_created:
        raise COULD_NOT_UPDATE_EXCEPTION(what="user's password")


@router.post('/testCreateUser')
async def test_create_user(
    username: str,
    email: str,
    password: str,
    response: Response,
    db: Session = Depends(get_db),
    csrf_protect: CsrfProtect = Depends()
) -> Token:
    hashed_password: str = hash_password(password)

    new_user: UserModel = await insert_user_into_db(
        username=username,
        email=email,
        hashed_password=hashed_password,
        db=db
    )
    if not new_user:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not create new user'
        )
    
    # Generate tokens
    access_token: str = generate_jwt_token(
        subject=str(new_user.user_id),
        expire_time=timedelta(minutes=ACCESS_TOKEN_LIFETIME_MINUTES),
        secret_key=ACCESS_TOKEN_SECRET
    )
    refresh_token: str = generate_jwt_token(
        subject=str(new_user.user_id),
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
