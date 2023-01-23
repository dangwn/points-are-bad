from fastapi import Request, HTTPException, status 
from typing import Optional
from fastapi.security.utils import get_authorization_scheme_param

from authentication.schema import HeaderCredentials

async def pab_auth_handler(request: Request) -> Optional[HeaderCredentials]:
    '''
    Retrieves the access token and provider from the request headers
    Pinched and adapted from the native ouath2scheme
    Expects the request to have headers:
    {
        "Authorization":"Bearer {token}",
        "provider":"{provider}"
    }
    '''
    authorization: str = request.headers.get('Authorization')
    provider: str = request.headers.get('provider')

    auth_scheme, auth = get_authorization_scheme_param(authorization)

    if auth_scheme.lower() != 'bearer' or not authorization or not provider:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated",
            headers={"WWW-Authenticate": "Bearer"},
        )

    return HeaderCredentials(access_token=auth, provider=provider)
    