from fastapi import (
    APIRouter,
    Depends,
    status,
    HTTPException
)

from auth.utils import get_current_user
from db import get_db
from points.schema import (
    UserWithPoints,
    LeaderBoardUser
)
from points.models import Points as PointsModel
from points.utils import (
    get_points_by_user_id,
    get_leaderboard
)
from user.models import User as UserModel

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
    user_points: Optional[PointsModel] = await get_points_by_user_id(
        user_id=current_user.user_id,
        db=db
    )
    
    if not user_points:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail='Could not get points for user'
        )
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
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail='Could not retrieve leaderboard'
        )
    return leaderboard
