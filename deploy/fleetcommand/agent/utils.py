import base64
import json
import os
import yaml
from pathlib import Path

def base64_encode(s: str) -> str:
    '''
    Base64 encodes a string
    '''
    return base64.b64encode(
        str.encode(s)
    ).decode('utf-8')

def create_docker_auth_secret(
    registry: str,
    username: str,
    password: str,
    email: str,
) -> str:
    '''
    Creates base64 encoded docker auth token to be used by k8s cluster
    '''
    token = f'{username}:{password}'
    encoded_token = base64_encode(token)

    auth_creds = {
        "auths":{
            registry:{
                "username":username,
                "password":password,
                "email":email,
                "auth": encoded_token
            }
        }
    }
    return base64_encode(
        json.dumps(
            auth_creds
        ))

def create_auth_secret_key() -> str:
    '''
    Creates JWT encoder secret key
    '''
    return os.popen('openssl rand -hex 32').read()

def create_values_yaml(
    pab_yaml_path: str = 'deploy/pab_master.yml',
    values_base_yaml_path: str = 'deploy/values_base.yml',
    values_yaml_path: str = 'deploy/helm/points-are-bad/values.yaml',
    write_yaml: bool = True
):
    pab_dict = yaml.safe_load(Path(pab_yaml_path).read_text())
    values_dict = yaml.safe_load(Path(values_base_yaml_path).read_text())

    values_dict['global']['appName'] = pab_dict['global']['appName']
    values_dict['backend']['api']['replicas'] = pab_dict['backend']['api']['replicas']
    values_dict['nginx']['nodePort'] = pab_dict['nginx']['nodePort']

    # Secrets
    values_dict['secrets']['authSecretKey'] = base64_encode(create_auth_secret_key())
    values_dict['secrets']['dbName'] = base64_encode(pab_dict['secrets']['dbName'])
    values_dict['secrets']['dbUser'] = base64_encode(pab_dict['secrets']['dbUser'])
    values_dict['secrets']['dbPassword'] = base64_encode(pab_dict['secrets']['dbPassword'])
    values_dict['secrets']['redisPassword'] = base64_encode(pab_dict['secrets']['redisPassword'])
    values_dict['secrets']['accessTokenSecret'] = base64_encode(create_auth_secret_key())
    values_dict['secrets']['refreshTokenSecret'] = base64_encode(create_auth_secret_key())
    values_dict['secrets']['csrfTokenSecret'] = base64_encode(create_auth_secret_key())

    # Docker
    values_dict['secrets']['dockerConfig'] = create_docker_auth_secret(
        registry = pab_dict['docker']['registry'],
        username = pab_dict['docker']['username'],
        password = pab_dict['docker']['password'],
        email = pab_dict['docker']['email']
    )

    if write_yaml:
        yaml.dump(
            values_dict,
            open(values_yaml_path,'w'),
            default_flow_style = False
        )
    return values_dict
