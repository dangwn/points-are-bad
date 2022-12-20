from fastapi import APIRouter, Depends, status
from sqlalchemy.orm import Session
from typing import List

from db import get_db
from authentication.utils import get_current_user, is_current_user_admin
from predictions import utils
from user.models import User as UserModel
from http_exceptions import NOT_AUTHORIZED_EXCEPTION, NOT_FOUND_EXCEPTION, NOT_ADMIN_EXCEPTION
from predictions.schema import MatchPrediction, MatchPredictionWithUserId, MatchPredictionWithTeams

router = APIRouter(
    prefix = '/predictions',
    tags = ['predictions']
)

@router.get('/', response_model = List[MatchPredictionWithTeams])
async def get_current_user_predictions(
    database: Session = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
) -> List[MatchPredictionWithTeams]:
    '''
    Returns the predictions for a given user
    '''

    predictions = await utils.get_display_predictions(current_user.id, database)

    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for user {current_user.id}')

    return predictions

@router.put('/', status_code = status.HTTP_202_ACCEPTED, response_model = List[MatchPredictionWithUserId])
async def update_predictions(
    predictions: List[MatchPredictionWithUserId],
    database = Depends(get_db),
    current_user: UserModel = Depends(get_current_user)
) -> List[MatchPredictionWithUserId]:
    '''
    Updates predictions(s) for a given user and match
    '''
    # Normal users can only update their own predictions
    # Admins can update all predictions
    users_from_predictions = list({pred.user_id for pred in predictions})
    if current_user is None:
        raise NOT_AUTHORIZED_EXCEPTION
    elif not current_user.is_admin:
        if not (len(users_from_predictions) == 1 and current_user.id in users_from_predictions):
            raise NOT_AUTHORIZED_EXCEPTION

    updated_predictions = await utils.update_predictions(predictions, database)
    return updated_predictions

@router.get('/user/{user_id}', response_model = List[MatchPrediction])
async def get_user_predictions(
    user_id: int,
    database: Session = Depends(get_db),
    is_admin: UserModel = Depends(is_current_user_admin)
) -> List[MatchPrediction]:
    '''
    Returns the predictions for a given user
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    predictions = await utils.get_predictions_by_user_id(user_id, database)

    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for user {user_id}')

    return predictions

@router.get('/match/{match_id}', response_model = List[MatchPredictionWithUserId])
async def get_match_predictions(
    match_id: int,
    database: Session = Depends(get_db),
    is_admin: bool = Depends(is_current_user_admin)
) -> List[MatchPredictionWithUserId]:
    '''
    Returns a list of predictions for a given match
    '''
    if not is_admin:
        raise NOT_ADMIN_EXCEPTION

    predictions = await utils.get_predictions_by_match_id(match_id, database)
    
    if predictions is None:
        raise NOT_FOUND_EXCEPTION(f'Predictions for match {match_id}')

    return predictions

