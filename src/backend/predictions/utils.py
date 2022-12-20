from typing import List
from sqlalchemy import and_
from sqlalchemy.orm import Session

from predictions.models import Prediction as PredictionModel
from predictions.schema import MatchPredictionWithUserId
from matches.models import Match as MatchModel
from http_exceptions import NOT_FOUND_EXCEPTION, COULD_NOT_UPDATE_EXCEPTION

async def get_predictions_by_user_id(user_id: int, database: Session) -> List[PredictionModel]:
    '''
    Queries database for predictions with given user id
    '''
    user_predictions = database.query(PredictionModel).filter(PredictionModel.user_id == user_id).all()

    if user_predictions is None:
        raise NOT_FOUND_EXCEPTION('Predictions for user {user_id}')

    return user_predictions

async def get_predictions_by_match_id(match_id: int, database: Session) -> List[PredictionModel]:
    '''
    Queries database for predictions with given match id
    '''
    match_predictions = database.query(PredictionModel).filter(PredictionModel.match_id == match_id).all()

    if match_predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for match {match_id}')

    return match_predictions

async def update_predictions(
    predictions: List[MatchPredictionWithUserId],
    database: Session
) -> List[PredictionModel]:
    '''
    Updates multiple predictions in the database
    '''
    updated_predictions = []
    try:    
        for pred in predictions:
            db_prediction = database.query(
                PredictionModel
            ).filter(
                and_(
                    PredictionModel.match_id == pred.match_id,
                    PredictionModel.user_id == pred.user_id,
                )
            ).first()

            if not db_prediction:
                continue
            
            db_prediction.predicted_home_goals = pred.predicted_home_goals
            db_prediction.predicted_away_goals = pred.predicted_away_goals

            updated_predictions.append(db_prediction)
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table with new predictions')
    
    database.commit()

    return updated_predictions

async def get_display_predictions(
    user_id: int,
    database: Session
):
    subq = database.query(
        PredictionModel
    ).filter(PredictionModel.user_id == user_id).subquery()

    predictions = database.query(
        subq,
        MatchModel.home,
        MatchModel.away
    ).join(
        MatchModel,
        subq.c.match_id == MatchModel.match_id,
        isouter = True
    ).all()

    return predictions