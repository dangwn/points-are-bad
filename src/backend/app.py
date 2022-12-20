from fastapi import FastAPI

from user.router import router as user_router
from authentication.router import router as login_router
from points.router import router as points_router
from matches.router import router as matches_router
from predictions.router import router as predictions_router

app = FastAPI()
app.include_router(user_router)
app.include_router(login_router)
app.include_router(points_router)
app.include_router(matches_router)
app.include_router(predictions_router)

@app.get('/')
async def root():
    return {'message':'Hello'}