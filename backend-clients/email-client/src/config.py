import os

CONSUMER_API_HOST: str = os.getenv('CONSUMER_API_HOST', '0.0.0.0')
CONSUMER_API_PORT: int = int(os.getenv('CONSUMER_API_PORT', 8002))

SMTP_USER_ACCESS_KEY_ID: str = os.getenv('AWS_ACCESS_KEY_ID')
SMTP_USER_SECRET_ACCESS_KEY_ID: str = os.getenv('AWS_SECRET_ACCESS_KEY')
AWS_REGION: str = os.getenv('AWS_REGION', 'eu-west-1')

EMAIL_CHARSET: str = 'UTF-8'
EMAIL_SENDER_ADDRESS: str = os.getenv('EMAIL_SENDER_ADDRESS')
EMAIL_SUBJECT: str = os.getenv('EMAIL_SUBJECT', 'Verify Points are Bad Email')

RABBITMQ_CONNECTION_STRING: str = os.getenv('RABBITMQ_CONNECTION_STRING', 'http://localhost:5672')
RABBITMQ_QUEUE_NAME: str = os.getenv('RABBITMQ_QUEUE_NAME', 'email_client_queue')