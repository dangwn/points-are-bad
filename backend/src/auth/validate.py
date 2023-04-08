from fastapi import HTTPException, status
from sqlalchemy.orm import Session

from auth.schema import LoginUser
from auth.utils import verify_password
from user.models import User as UserModel


from typing import (
    Optional
)

async def validate_login_user(
    login_user: LoginUser,
    db: Session
) -> Optional[UserModel]:
    '''
    Function called by login endpoint to validate user's email and password
    '''
    user: Optional[UserModel] = db.query(UserModel).filter(
        UserModel.email == login_user.email
    ).first()
    if not user:
        return
    
    if not verify_password(
        plain_password=login_user.password,
        hashed_password=user.hashed_password
    ):
        return
    return user