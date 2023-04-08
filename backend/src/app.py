from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from config import REQUEST_ORIGINS

from auth.router import router as AuthRouter
from match.router import router as MatchRouter
from points.router import router as PointsRouter
from user.router import router as UserRouter

from typing import Dict

app: FastAPI = FastAPI()

app.include_router(router=AuthRouter)
app.include_router(router=MatchRouter)
app.include_router(router=PointsRouter)
app.include_router(router=UserRouter)

app.add_middleware(
    CORSMiddleware,
    allow_origins=REQUEST_ORIGINS,
    allow_credentials=True,
    allow_methods=['*'],
    allow_headers=['*']
)

@app.get('/')
async def root() -> Dict[str,str]:
    return {'message':'We are go'}
