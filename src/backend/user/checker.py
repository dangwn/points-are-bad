from user.models import User
from sqlalchemy.orm import Session

async def check_if_email_exists(email: str, database: Session):
    '''
    Check to see if any users in db have the given email
    '''
    return database.query(User).filter(User.email == email).first()


async def check_if_username_exists(username: str, database: Session):
    '''
    Check to see if any users in db have given username
    '''
    return database.query(User).filter(User.username == username).first()
