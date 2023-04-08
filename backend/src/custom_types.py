from typing import Type, NewType
from sqlalchemy.ext.declarative import DeclarativeMeta

DeclarativeBase = NewType('DeclarativeBase', Type[DeclarativeMeta])