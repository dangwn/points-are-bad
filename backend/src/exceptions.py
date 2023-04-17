from fastapi import HTTPException, status

COULD_NOT_CREATE_MATCH_EXCEPTION: HTTPException = HTTPException(
    status_code=status.HTTP_400_BAD_REQUEST,
    detail='Could not create match.'
)

NO_CURRENT_USER_EXCEPTION: HTTPException = HTTPException(
    status_code=status.HTTP_403_FORBIDDEN,
    detail='Could not retrieve current user.'
)

USER_IS_NOT_ADMIN_EXCEPTION: HTTPException = HTTPException(
    status_code=status.HTTP_401_UNAUTHORIZED,
    detail='User is not an admin.'
)

def COULD_NOT_UPDATE_EXCEPTION(what: str, why: str) -> HTTPException:
    return HTTPException(
        status_code=status.HTTP_406_NOT_ACCEPTABLE,
        detail=f'Could not update {what} due to {why}.'
    )