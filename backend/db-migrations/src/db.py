from sqlalchemy import create_engine
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

from config import (
    DB_USER,
    DB_PASSWORD,
    DB_HOST,
    DB_PORT,
    DB_NAME,
    DB_TYPE
)

from sqlalchemy.engine import Engine
from custom_types import DeclarativeBase

SQLALCHEMY_DATABASE_URL: str = f"{DB_TYPE}://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}"

engine: Engine = create_engine(SQLALCHEMY_DATABASE_URL)

SessionLocal: sessionmaker = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base: DeclarativeBase = declarative_base()

def get_db() -> None:
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()