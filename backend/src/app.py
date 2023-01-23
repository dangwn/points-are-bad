from fastapi import FastAPI

from user.router import router as user_router
from points.router import router as points_router
from matches.router import router as matches_router
from predictions.router import router as predictions_router

app = FastAPI()

app.include_router(user_router)
app.include_router(points_router)
app.include_router(matches_router)
app.include_router(predictions_router)

from fastapi.middleware.cors import CORSMiddleware
origins = ['http://localhost','http://localhost:3000']
app.add_middleware(
    CORSMiddleware,
    allow_origins = origins,
    allow_credentials = True,
    allow_methods = ['*'],
    allow_headers = ['*']
)

@app.get('/')
async def root():
    return {'message':'Hello'}
