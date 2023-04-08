import re
from sqlalchemy.orm import Session

from user.models import User as UserModel
from config import (
    USERNAME_MIN_LENGTH,
    USERNAME_MAX_LENGTH
)

def validate_email(email: str) -> bool:
    '''
    Validates email format
    '''  
    if re.match(r"[^@]+@[^@]+\.[^@]+", email):  
        return True  
    return False   

async def validate_user_entries(
    username: str,
    email: str
) -> bool:
    '''
    Validates user email and username
    @TODO: Add username filtering
    '''
    if not validate_email(email):
        return False    
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