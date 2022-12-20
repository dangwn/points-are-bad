from typing import List, Optional
from sqlalchemy.orm import Session

from user.models import User as UserModel
from matches.models import Match as MatchModel
from predictions.models import Prediction as PredictionModel
from http_exceptions import COULD_NOT_UPDATE_EXCEPTION

async def populate_predictions(
    database: Session,
    user_id: Optional[int] = None,
    match_id: Optional[int] = None
) -> Optional[List[PredictionModel]]:
    '''
    Populates the predictions table for a given user or match
    '''
    if user_id is None and match_id is None:
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table')

    model = UserModel if user_id is None else MatchModel

    items = database.query(model).all()
    if items == []:
        return
    
    new_predictions = []
    for item in items:
        if user_id is not None:
            new_pred = PredictionModel(user_id=user_id, match_id=item.match_id)
        else:
            new_pred = PredictionModel(user_id=item.id, match_id=match_id)
        try:
            database.add(new_pred)

            new_predictions.append(new_pred)
        except:
            pass

    database.commit()
    if new_predictions != []:
        [database.refresh(pred) for pred in new_predictions]
    
    return new_predictions

    