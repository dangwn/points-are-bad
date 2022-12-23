from urllib import request, parse
from typing import Optional
import json

host = 'http://localhost:8000/'
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
    data: Optional[dict] = None,
    token: Optional[str] = None,
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

def login(
    username: str = 'dan',
    password: str = 'password'
):
    encoded_data = parse.urlencode({
        'grant_type': '',
        'username': username,
        'password': password,
        'scope': '',
        'client_id': '',
        'client_secret': '',
    }).encode()
    headers = {
        'accept': 'application/json',
        'Content-Type': 'application/x-www-form-urlencoded'
    }
    req = request.Request(
        url = host + 'login/',
        data = encoded_data,
        headers = headers
    )
    response = request.urlopen(req)
    return response.status, json.loads(response.read().decode('utf-8'))['access_token']

def create_user(
    username: str = 'dan',
    email: str = 'dan@email.com',
    password: str = 'password'
):
    return url_request(
        'user/',
        data = {
            'username': username,
            'email': email,
            'password': password
        }
    )

def create_match(
    match_date: str = '2023-01-01',
    home: str = 'Eng',
    away: str = 'Ger',
    token: str = ''
):
    return url_request(
        endpoint = 'matches/', 
        data = {
            'match_date': match_date,
            'home': home,
            'away': away
        },
        token = token)

def update_score(
    user_id: int = 1,
    points: int = 10,
    correct_scores: int = 2,
    largest_error: int = 5,
    token: str = ''
):
    return url_request(
        f'points/player/{user_id}',
        data = {
            'points': points,
            'correct_scores': correct_scores,
            'largest_error': largest_error
        },
        token = token,
        method = 'PUT'
    )

def get_position(
    token: str
):
    return url_request(
        'points/position',
        token = token,
        method = 'GET'
    )

def get_leaderboard(
    table_start: int = 1,
    table_end: int = 100
):
    pass

def main():
    print(create_user())
    print(create_user('string','string@example.com','string'))

    _, token = login('dan', 'password')
    print(create_match(token = token))
    print(create_match(token = token))
    print(create_match(token = token))

    print(update_score(token = token))
    print(get_position(token = token))

if __name__ == '__main__':
    main()

