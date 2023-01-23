from urllib import request, parse
from typing import Optional
import json

host = 'http://localhost:30009/'
default_headers = {
    'Content-Type':'application/json',
    'accept':'application/json'
}

def encode_data(data: dict):
    end_of_items = len(data) - 1
    encoded_str = '{'

    for ind, (key, value) in enumerate(data.items()):
        if type(value) == int or type(value) == float:
            str_value = value
        elif type(value) == bool:
            str_value = str(value).lower()
        else:
            str_value = f'"{value}"'
        encoded_str += f'"{key}":{str_value}'

        if ind != end_of_items:
            encoded_str += ', '

    encoded_str += '}'
    return encoded_str.encode()
    
def url_request(
    endpoint: str,
    token: str,
    provider: str,
    data: Optional[dict] = None,
    host: str = host,
    method: str = 'POST'
):
    encoded_data = encode_data(data) if data is not None else None
    headers = default_headers.copy()
    if token:
        headers['Authorization'] = f'Bearer {token}'

    req = request.Request(
        url = host + endpoint,
        data = encoded_data,
        headers = headers,
        method = method
    )

    response = request.urlopen(req)
    return response.status, json.loads(response.read().decode('utf-8'))

def create_user(

)