import glob
import logging
import os
from typing import List, Optional

import alembic.config

logging.basicConfig(level=logging.DEBUG)
    
def get_init_dir() -> str:
    '''
    Finds the alembic init file
    '''
    init_path: List[str] = glob.glob(os.path.join('**','alembic.ini'), recursive=True)
    if init_path == []:
        raise FileNotFoundError('Could not find alembic init file')
    if len(init_path) > 1:
        raise FileNotFoundError(f'Found multiple alembic init files: {init_path}')

    return os.path.dirname(init_path[0])

def run_alembic(
    alembic_args: List[str],
    init_dir: str = '',
    log_msg: Optional[str] = None
) -> None:
    '''
    Runs an alembic command against a given init file
    '''
    old_path: str = os.getcwd()
    if init_dir != '': os.chdir(init_dir)

    if log_msg is not None:
        logging.info(log_msg)
    alembic.config.main(argv=alembic_args)

    os.chdir(old_path)

def main() -> None:    
    init_dir: Optional[str] = get_init_dir()
    try:
        os.mkdir(os.path.join(
            os.getcwd(),
            init_dir,
            'alembic',
            'versions'
        ))
    except:
        pass

    run_alembic(['stamp','head'], init_dir, 'Stamping head')
    run_alembic(['revision','--autogenerate'], init_dir, 'Creating new db revision')
    run_alembic(['upgrade','head'], init_dir, 'Upgrading db')

if __name__ == '__main__':
    main()