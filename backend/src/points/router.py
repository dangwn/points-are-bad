from fastapi import (
    APIRouter,
    Depends,
    status,
    Response
)

from auth.utils import get_current_user, is_current_user_admin
from db import get_db
from points.schema import (
    UserWithPoints,
    LeaderBoardUser
)
from points.models import Points as PointsModel
from points.utils import (
    get_points_by_user_id,
    get_leaderboard,
    update_all_points
)
from user.models import User as UserModel
from exceptions import (
    COULD_NOT_GET_POINTS_EXCEPTION,
    COULD_NOT_GET_LEADERBOARD_EXCEPTION,
    NO_CURRENT_USER_EXCEPTION,
    USER_IS_NOT_ADMIN_EXCEPTION,
    COULD_NOT_UPDATE_EXCEPTION
)

from typing import (
    Optional,
    List
)
from sqlalchemy.orm import Session

router: APIRouter = APIRouter(
    prefix='/points',
)

@router.get('/', response_model=UserWithPoints)
async def get_user_points(
    current_user: Optional[UserModel] = Depends(get_current_user),
    db: Session = Depends(get_db)
) -> UserWithPoints:
    if not current_user:
        raise NO_CURRENT_USER_EXCEPTION
    
    user_points: Optional[PointsModel] = await get_points_by_user_id(
        user_id=current_user.user_id,
        db=db
    )
    
    if not user_points:
        raise COULD_NOT_GET_POINTS_EXCEPTION
    return user_points

@router.get('/leaderboard', response_model=List[LeaderBoardUser])
async def get_global_leaderboard(
    limit: int = 10,
    offset: int = 0,
    db: Session = Depends(get_db)
) -> List[LeaderBoardUser]:
    leaderboard: List[PointsModel] = await get_leaderboard(
        limit=limit,
        offset=offset,
        db=db
    )

    if not leaderboard:
        raise COULD_NOT_GET_LEADERBOARD_EXCEPTION
    return leaderboard

@router.post('/calculate', status_code=status.HTTP_204_NO_CONTENT)
async def calculate_points(
    user_is_admin: bool = Depends(is_current_user_admin),
    db: Session = Depends(get_db)
) -> Response:
    if not user_is_admin:
        raise USER_IS_NOT_ADMIN_EXCEPTION
    
    updated: bool = await update_all_points(db)
    if not updated:
        raise COULD_NOT_UPDATE_EXCEPTION(what='points table', why="I'm bad at sql")