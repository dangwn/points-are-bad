import requests
from requests import Response
import os
from typing import Optional, Dict, List, Any

from google.oauth2 import id_token
from google.auth.transport import requests as google_requests

from config import GOOGLE_CLIENT_ID

async def validate_github_token(token: str) -> Optional[str]:
    '''
    Validates the github token and returns the email if valid
    '''
    try:
        headers: Dict[str, str] = {
            'Authorization': f'Bearer {token}',
        }

        response: Response = requests.get('https://api.github.com/user/emails', headers=headers)
        email_info: List[Dict[str, str]] = response.json()

        for info in email_info:
            if info['primary'] == True and info['verified'] == True:
                return info['email']

        raise Exception("Could not find github user's primary, verified email")
    except Exception as e:
        print(f'Error when validating github token: {e}')
        return

async def validate_google_token(token: str) -> Optional[str]:
    '''
    Validates the google token and returns the email if valid
    '''    
    try:
        response: Dict[str, Any] = id_token.verify_oauth2_token(
            id_token=token,
            request=google_requests.Request(),
            audience=GOOGLE_CLIENT_ID
        )
        return response['email']
    except Exception as e:
        print(f'Error when validating google token: {e}')
        return

async def validate_test_user(token: str):
    return 'test@test.com'

async def validate_token(token: str, provider: str) -> Optional[str]:
    if provider == 'google':
        google_email: Optional[str] = await validate_google_token(token)
        return google_email
    elif provider == 'github':
        github_email: Optional[str] = await validate_github_token(token)
        return github_email
    elif provider == 'test':
        test_email: Optional[str] = await validate_test_user(token)
        return test_email

    return
