import asyncio

from .custom_types import Callback
from typing import Callable, Coroutine, List, Union

def is_async(func: Union[Callable, Coroutine]) -> bool:
    return asyncio.iscoroutinefunction(func)

async def run_callback(callback: Callback) -> None:
    if is_async(callback):
        await callback()
    else:
        callback()

async def run_callbacks(callbacks: List[Callback]) -> None:
    for cb in callbacks:
        await run_callback(cb)