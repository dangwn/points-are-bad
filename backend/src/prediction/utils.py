from prediction.models import Prediction as PredictionModel
from match.models import Match as MatchModel

from datetime import date
from sqlalchemy.orm import Session, Query
from typing import Optional, List

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