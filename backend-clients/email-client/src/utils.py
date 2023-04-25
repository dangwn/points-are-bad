import boto3
import json
import re

from config import (
    AWS_REGION,
    SMTP_USER_ACCESS_KEY_ID,
    SMTP_USER_SECRET_ACCESS_KEY_ID,
    EMAIL_CHARSET,
    EMAIL_SENDER_ADDRESS,
    EMAIL_SUBJECT
)

from aio_pika.message import IncomingMessage
from botocore.client import BaseClient
from botocore.exceptions import ClientError
from typing import Dict, Union
from re import Pattern

client: BaseClient = boto3.client(
    'ses', 
    region_name=AWS_REGION, 
    aws_access_key_id=SMTP_USER_ACCESS_KEY_ID,
    aws_secret_access_key=SMTP_USER_SECRET_ACCESS_KEY_ID
)

def clean_html(raw_html) -> str:
    cleaner: Pattern = re.compile('<.*?>')
    clean_text: str = re.sub(cleaner, '', raw_html)
    return clean_text

def create_email_content(verification_link: str) -> str:
    return f'''<html>
<head></head>
<body>
  <h1>Please Verify Your Email</h1>
  <p>Please follow the link below to verify your email:</p>
  <p>{verification_link}</p>
</body>
</html>
    '''

def send_email(address: str, verification_link: str) -> None:
    email_html: str = create_email_content(verification_link)

    try:
        response = client.send_email(
            Destination={
                'ToAddresses': [
                    address,
                ],
            },
            Message={
                'Body': {
                    'Html': {
                        'Charset': EMAIL_CHARSET,
                        'Data': email_html,
                    },
                    'Text': {
                        'Charset': EMAIL_CHARSET,
                        'Data': clean_html(email_html),
                    },
                },
                'Subject': {
                    'Charset': EMAIL_CHARSET,
                    'Data': EMAIL_SUBJECT,
                },
            },
            Source=EMAIL_SENDER_ADDRESS,
        )
    except ClientError as e:
        print(e.response['Error']['Message'])
    else:
        print("Email sent! Message ID:"),
        print(response['MessageId'])

async def email_callback(message: IncomingMessage) -> None:
    await message.ack()
    body: Union[Dict[str,str], str] = json.loads(message.body.decode())
    if type(body) != dict:
        print('Could not decode body')
        return
    
    email_address: str = next(iter(body.keys()))
    verification_link: str = body[email_address]

    send_email(
        address=email_address,
        verification_link=verification_link
    )

    

    