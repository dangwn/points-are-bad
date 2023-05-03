import asyncio

from typing import Callable, List

def is_async(func: Callable) -> bool:
    return asyncio.iscoroutinefunction(func)

async def run_callback(callback: Callable[..., None]) -> None:
    if is_async(callback):
        await callback()
    else:
        callback()

async def run_callbacks(callbacks: List[Callable[..., None]]) -> None:
    for cb in callbacks:
        await run_callback(cb)