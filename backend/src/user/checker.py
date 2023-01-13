from user.models import User
from sqlalchemy.orm import Session

async def email_exists(email: str, database: Session):
    '''
    Check to see if any users in db have the given email
    '''
    email = database.query(User).filter(User.email == email).first()
    if email:
        return True
    return False


async def username_exists(username: str, database: Session):
    '''
    Check to see if any users in db have given username
    '''
    username = database.query(User).filter(User.username == username).first()
    if username:
        return True
    return False

async def is_first_user(database: Session):
    first_user = database.query(User).first() 
    if first_user:
        return False
    return True
     