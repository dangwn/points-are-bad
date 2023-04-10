import re
from sqlalchemy.orm import Session

from user.models import User as UserModel
from config import (
    USERNAME_MIN_LENGTH,
    USERNAME_MAX_LENGTH
)

async def validate_user_entries(
    username: str,
    email: str
) -> bool:
    '''
    Validates user email and username
    @TODO: Add username filtering
    '''
    if len(username) < USERNAME_MIN_LENGTH or len(username) > USERNAME_MAX_LENGTH:
        return False
    return True

async def is_user_email_in_db(
    email: str,
    db: Session
) -> bool:
    '''
    Checks whether there is a user in the database that has the
        provided email
    '''
    if db.query(UserModel).filter(
        UserModel.email == email
    ).first():
        return True
    return False