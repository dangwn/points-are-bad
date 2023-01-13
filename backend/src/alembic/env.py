'''
This template is updated following this tutorial: 
https://www.jetbrains.com/pycharm/guide/tutorials/fastapi-aws-kubernetes/setup_sqlalchemy_2/
'''

from logging.config import fileConfig

from sqlalchemy import engine_from_config
from sqlalchemy import pool

from alembic import context

config = context.config

if config.config_file_name is not None:
    fileConfig(config.config_file_name)

# Points to the backend
from db import Base
import config as pab_config
from user.models import User
from points.models import PlayerPoints
from matches.models import Match
from predictions.models import Prediction

target_metadata = Base.metadata

def get_db_url():
    DB_USER = pab_config.DB_USER
    DB_PASSWORD = pab_config.DB_PASSWORD
    DB_HOST = pab_config.DB_HOST
    DB_PORT = pab_config.DB_PORT
    DB_NAME = pab_config.DB_NAME
    DB_TYPE = pab_config.DB_TYPE

    return f"{DB_TYPE}://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"

def run_migrations_offline() -> None:
    """Run migrations in 'offline' mode.

    This configures the context with just a URL
    and not an Engine, though an Engine is acceptable
    here as well.  By skipping the Engine creation
    we don't even need a DBAPI to be available.

    Calls to context.execute() here emit the given string to the
    script output.

    """
    url = get_db_url() # Update with own url
    context.configure(
        url=url,
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )

    with context.begin_transaction():
        context.run_migrations()


def run_migrations_online() -> None:
    """Run migrations in 'online' mode.

    In this scenario we need to create an Engine
    and associate a connection with the context.

    """

    # This block is taken from url at top
    configuration = config.get_section(config.config_ini_section)
    configuration["sqlalchemy.url"] = get_db_url()
    connectable = engine_from_config(
        configuration, prefix="sqlalchemy.", poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(
            connection=connection, target_metadata=target_metadata
        )

        with context.begin_transaction():
            context.run_migrations()


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()
