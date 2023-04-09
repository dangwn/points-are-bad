import requests

HOST = 'http://localhost:8000'
HEADERS = {
    'accept': 'application/json',
    'Content-Type': 'application/json',
}

def add_user(
    email: str,
    username: str,
    password: str
):

    json_data = {
        'username': username,
        'email': email,
        'password': password,
    }
    r = requests.post(
        url=f'{HOST}/user/',
        headers=HEADERS,
        json=json_data
    )
    return r.status_code
    
def create_match(
    match_date: str,
    home: str,
    away: str
):
    json_data = {
        'match_date': match_date,
        'home': home,
        'away': away,
    }

    r = requests.post(
        url=f'{HOST}/match/',
        headers=HEADERS,
        json=json_data
    )
    return r.status_code

requests.get(HOST)
add_user('dan@dan.com','dan','password')
add_user('chris@chris.com','chris','password')
add_user('b@b.com','bbb','password')

create_match('2023-01-01', home='Eng',away='Ger')
create_match('2023-06-01', home='Fra',away='Bra')
create_match('2023-05-03', home='Spa',away='Ita')