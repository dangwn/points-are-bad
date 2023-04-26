from typing import Literal, Callable, Coroutine, Union, Any, TypeVar

Callback = Union[Callable[..., None], Coroutine[Any, Any, None]]

DecoratedCallable = TypeVar('DecoratedCallable', bound=Callable[..., Any])

EventType = Literal['startup','loop','shutdown']
