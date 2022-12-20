from fastapi import HTTPException, status

CREDENTIALS_EXCEPTION = HTTPException(
    status_code = status.HTTP_401_UNAUTHORIZED,
    detail = 'Could not validate user token',
    headers = {'WWW-Authenticate':'Bearer'}
)

EMAIL_EXISTS_EXCEPTION = HTTPException(
    status_code = status.HTTP_400_BAD_REQUEST,
    detail = f'User with given username already exists'
)

MATCH_NOT_FOUND_EXCEPTION = HTTPException(
    status_code = status.HTTP_404_NOT_FOUND,
    detail = 'Match not found'
)

NOT_ADMIN_EXCEPTION = HTTPException(
    status_code = status.HTTP_401_UNAUTHORIZED,
    detail = 'Not admin user'
)

NOT_AUTHORIZED_EXCEPTION = HTTPException(
    status_code = status.HTTP_401_UNAUTHORIZED,
    detail = 'User is not authorized'
)

PASSWORD_INCORRECT_EXCEPTION = HTTPException(
    status_code = status.HTTP_403_FORBIDDEN,
    detail = 'Password incorrect'
)

USERNAME_EXISTS_EXCEPTION = HTTPException(
    status_code = status.HTTP_400_BAD_REQUEST,
    detail = f'User with given email already exists'
)

USER_NOT_FOUND_EXCEPION =  HTTPException(
    status_code = status.HTTP_404_NOT_FOUND,
    detail = 'User not found'
)

USERNAME_OR_PASSWORD_EXCEPTION = HTTPException(
    status_code = status.HTTP_401_UNAUTHORIZED,
    detail = 'Incorrect username or password',
    headers = {'WWW-Authenticate':'Bearer'}
)

def NOT_FOUND_EXCEPTION(what: str) -> HTTPException:
    return HTTPException(
        status_code = status.HTTP_404_NOT_FOUND,
        detail = f'{what} could not be found'
    )

def COULD_NOT_UPDATE_EXCEPTION(what: str) -> HTTPException:
    return HTTPException(
        status_code = status.HTTP_400_BAD_REQUEST,
        detail = f'Could not update {what}'
    )

def INVALID_QUERY_PARAMETERS_EXCEPTION(msg: str) -> HTTPException:
    return HTTPException(
        status_code = status.HTTP_400_BAD_REQUEST,
        detail = msg
    )