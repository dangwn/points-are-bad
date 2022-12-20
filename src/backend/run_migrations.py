import os
import alembic.config
import time, logging
from typing import List, Optional

from db import get_db
import glob

logging.basicConfig(level = logging.DEBUG)

def poll_db() -> None:
    '''
    Polls the database until a connection can be established
    '''
    logging.info('Attempting to verify database connection')
    while True:
        if get_db() is not None:
            break
        logging.info('Could not connect to database. Retrying in 5 seconds.')
        time.sleep(5)
    
def get_init_dir() -> str:
    '''
    Finds the alembic init file
    '''
    init_path = glob.glob(os.path.join('**','alembic.ini'), recursive=True)
    if init_path == []:
        raise FileNotFoundError('Could not find alembic init file')
    if len(init_path) > 1:
        raise FileNotFoundError(f'Found multiple alembic init files: {init_path}')

    return os.path.dirname(init_path[0])


def run_alembic(alembic_args: List[str], init_dir: str = '', log_msg: Optional[str] = None) -> None:
    '''
    Runs an alembic command against a given init file
    '''
    old_path = os.getcwd()
    if init_dir != '': os.chdir(init_dir)

    if log_msg is not None:
        logging.info(log_msg)
    alembic.config.main(argv = alembic_args)

    os.chdir(old_path)

def main():
    poll_db()
    
    init_dir = get_init_dir()

    run_alembic(['stamp','head'], init_dir, 'Stamping head')
    run_alembic(['revision','--autogenerate'], init_dir, 'Creating new db revision')
    run_alembic(['upgrade','head'], init_dir, 'Upgrading db')

if __name__ == '__main__':
    main()