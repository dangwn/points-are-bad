from prediction.models import Prediction as PredictionModel
from prediction.schema import Prediction
from match.models import Match as MatchModel

from datetime import date
from sqlalchemy.sql import case
from sqlalchemy.orm import Session, Query
from typing import (
    Dict,
    Optional,
    List,
    Union,
    Tuple, 
    Set
)

async def get_user_predictions_by_id(
    user_id: int,
    start_date: Optional[date],
    end_date: Optional[date],
    db: Session,
) -> List[PredictionModel]:
    predictions_query: Query = db.query(
        PredictionModel
    ).filter(
        PredictionModel.user_id == user_id
    ).join(MatchModel)
    
    if start_date:
        predictions_query = predictions_query.filter(MatchModel.match_date >= start_date)
    if end_date:
        predictions_query = predictions_query.filter(MatchModel.match_date <= end_date)

    predictions: List[PredictionModel] = predictions_query.order_by(
        MatchModel.match_date
    ).all()

    return predictions

async def update_user_predictions_by_id(
    user_id: int,
    new_predictions: List[Prediction],
    db: Session
) -> bool:
    home_goals_payload: Dict[int, Union[int, None]] = {
        pred.prediction_id: pred.home_goals for pred in new_predictions
    }
    away_goals_payload: Dict[int, Union[int, None]] = {
        pred.prediction_id: pred.away_goals for pred in new_predictions
    }

    # Check to see if the predictions to update are all from the current user
    predictions_to_update: List[Tuple[int]] = db.query(
        PredictionModel.user_id
    ).filter(
        PredictionModel.prediction_id.in_(home_goals_payload)
    ).all()
    if any(pred[0] != user_id for pred in predictions_to_update):
        return False

    try:
        db.query(
            PredictionModel
        ).filter(
            PredictionModel.prediction_id.in_(home_goals_payload)
        ).update({
            PredictionModel.home_goals: case(
                home_goals_payload,
                value=PredictionModel.prediction_id
            ),
            PredictionModel.away_goals: case(
                away_goals_payload,
                value=PredictionModel.prediction_id
            )
        }, synchronize_session=False)

        db.commit()
    except Exception as e:
        print(e)
        return False
    
    return True