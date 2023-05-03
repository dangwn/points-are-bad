from typing import Callable, TypeVar
from typing_extensions import Literal

DecoratedCallback = TypeVar('DecoratedCallback', bound=Callable[..., None])

EventType = Literal['startup','loop','shutdown']