from fastapi import (
    APIRouter,
    Depends,
    status,
    Response
)

from auth.utils import get_current_user
from db import get_db
from exceptions import (
    COULD_NOT_UPDATE_EXCEPTION,
    NO_CURRENT_USER_EXCEPTION
)
from prediction.models import Prediction as PredictionModel
from prediction.schema import (
    Prediction,
    PredictionWithMatch
)
from prediction.utils import (
    get_user_predictions_by_id,
    update_user_predictions_by_id
)
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

@router.put('/', status_code=status.HTTP_202_ACCEPTED)
async def update_user_predictions(
    new_predictions: List[Prediction],
    current_user: Optional[UserModel] = Depends(get_current_user),
    db: Session = Depends(get_db)
) -> Response:
    if new_predictions == []:
        return Response
    
    update_was_successful: bool = await update_user_predictions_by_id(
        user_id=current_user.user_id,
        new_predictions=new_predictions,
        db=db
    )
    if not update_was_successful:
        raise COULD_NOT_UPDATE_EXCEPTION(
            what='predictions table',
            why='erroneous predictions'
        )