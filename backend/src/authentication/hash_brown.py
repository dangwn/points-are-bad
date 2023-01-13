from passlib.context import CryptContext

pwd_context = CryptContext(schemes = ['bcrypt'], deprecated = 'auto')

def get_password_hash(password: str) -> str:
    '''
    Creates a hash for a given password
    '''
    return pwd_context.hash(password)

def verify_password(plain_password:str, hashed_password: str) -> bool:
    '''
    Verifies a plain text password against a hashed password
    '''
    return pwd_context.verify(plain_password, hashed_password)