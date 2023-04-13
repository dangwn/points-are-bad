from fastapi import HTTPException, status

NO_CURRENT_USER_EXCEPTION: HTTPException = HTTPException(
    status_code=status.HTTP_403_FORBIDDEN,
    detail='Could not retrieve current user.'
)

def COULD_NOT_UPDATE_EXCEPTION(what: str, why: str) -> HTTPException:
    return HTTPException(
        status_code=status.HTTP_406_NOT_ACCEPTABLE,
        detail=f'Could not update {what} due to {why}.'
    )