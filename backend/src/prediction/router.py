from fastapi import (
    APIRouter,
    Depends
)

from auth.utils import get_current_user
from db import get_db
from exceptions import NO_CURRENT_USER_EXCEPTION
from prediction.models import Prediction as PredictionModel
from prediction.schema import PredictionWithMatch
from prediction.utils import get_user_predictions_by_id
from user.models import User as UserModel

from datetime import date
from typing import (
    Optional,
    List
)
from sqlalchemy.orm import Session

router: APIRouter = APIRouter(
    prefix='/prediction'
)

@router.get('/', response_model=List[PredictionWithMatch])
async def get_user_predictions(
    current_user: Optional[UserModel] = Depends(get_current_user),
    start_date: Optional[date] = None,
    end_date: Optional[date] = None,
    db: Session = Depends(get_db)
) -> List[PredictionWithMatch]:
    if not current_user:
        raise NO_CURRENT_USER_EXCEPTION
    
    predictions: List[PredictionModel] = await get_user_predictions_by_id(
        user_id=current_user.user_id,
        start_date=start_date,
        end_date=end_date,
        db=db
    )
    
    return predictions