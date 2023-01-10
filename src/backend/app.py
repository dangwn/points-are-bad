from fastapi import FastAPI, Request
from fastapi.responses import JSONResponse

from user.router import router as user_router
from authentication.router import router as login_router
from points.router import router as points_router
from matches.router import router as matches_router
from predictions.router import router as predictions_router

from fastapi_jwt_auth import AuthJWT
from fastapi_jwt_auth.exceptions import AuthJWTException
from datetime import timedelta
from pydantic import BaseModel
import config

app = FastAPI()

class Settings(BaseModel):
    authjwt_secret_key: str = config.AUTH_SECRET_KEY
    authjwt_algorithm: str = config.AUTH_ALGORITHM
    authjwt_access_token_expires: timedelta = timedelta(minutes = config.AUTH_ACCESS_TOKEN_EXPIRE_MINUTES)

    authjwt_refresh_token_expires: timedelta = timedelta(days = 30)
    authjwt_token_location: set = {'headers', 'cookies'}

@AuthJWT.load_config
def get_config():
    return Settings()

@app.exception_handler(AuthJWTException)
def authjwt_exception_handler(request: Request, exc: AuthJWTException):
    return JSONResponse(
        status_code=exc.status_code,
        content={'detail': exc.message}
    )


app.include_router(user_router)
app.include_router(login_router)
app.include_router(points_router)
app.include_router(matches_router)
app.include_router(predictions_router)


from fastapi.middleware.cors import CORSMiddleware
origins = ['http://localhost','http://localhost:5173']
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