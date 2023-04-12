from fastapi import HTTPException, status

NO_CURRENT_USER_EXCEPTION: HTTPException = HTTPException(
    status_code=status.HTTP_403_FORBIDDEN,
    detail='Could not retrieve current user'
)