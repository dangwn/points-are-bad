from typing import List
from sqlalchemy import and_
from sqlalchemy.orm import Session
from sqlalchemy.sql import case
from datetime import datetime, date

from predictions.models import Prediction as PredictionModel
from predictions.schema import MatchPredictionWithUserId
from matches.models import Match as MatchModel
from http_exceptions import NOT_FOUND_EXCEPTION, COULD_NOT_UPDATE_EXCEPTION, NOTHING_TO_UPDATE_EXCEPTION

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
    user_id: int,
    database: Session
) -> List[PredictionModel]:
    '''
    Updates multiple predictions in the database
    '''
    home_goals_payload = {pred.match_id:pred.predicted_home_goals for pred in predictions}
    away_goals_payload = {pred.match_id:pred.predicted_away_goals for pred in predictions}
    
    # If any predictions are for matches today or in the past, raise an error
    today = date.today()
    prediction_match_dates = database.query(
        PredictionModel.match_date
    ).filter(and_(
        PredictionModel.user_id == user_id,
        PredictionModel.match_id.in_(home_goals_payload)
    )).all()
    prediction_match_dates = {md[0] for md in prediction_match_dates}
    if not prediction_match_dates:
        raise NOTHING_TO_UPDATE_EXCEPTION

    
    if any(datetime.strptime(match_date, '%Y-%m-%d').date() <= today for match_date in prediction_match_dates):
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table as some predictions are for games in the past')

    # Update scores
    try:
        database.query(
            PredictionModel
        ).filter(and_(
            PredictionModel.user_id == user_id,
            PredictionModel.match_id.in_(home_goals_payload)
        )).update({
            PredictionModel.predicted_home_goals: case(
                home_goals_payload,
                value = PredictionModel.match_id
            ),
            PredictionModel.predicted_away_goals: case(
                away_goals_payload,
                value = PredictionModel.match_id
            ),
        }, synchronize_session = False)
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table')

    database.commit()

    return await get_display_predictions(user_id, database)

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
        MatchModel.away,
        MatchModel.match_date
    ).join(
        MatchModel,
        subq.c.match_id == MatchModel.match_id,
        isouter = True
    ).all()

    return predictions