from typing import List
from sqlalchemy import func
from sqlalchemy.orm import Session
from sqlalchemy.sql import case

from http_exceptions import COULD_NOT_UPDATE_EXCEPTION, MATCH_NOT_FOUND_EXCEPTION
from matches.models import Match as MatchModel
from matches import schema
from predictions.populate import populate_predictions

async def get_all_matches(database: Session) -> List[MatchModel]:
    '''
    Gets all matches from database
    '''
    all_matches = database.query(MatchModel).all()
    return all_matches

async def create_new_match(match: schema.Match, database: Session) -> MatchModel:
    '''
    Creates new match and inserts it into database
    '''
    new_match = MatchModel(
        match_date = match.match_date,
        home = match.home,
        away = match.away
    )

    try:
        database.add(new_match)
        database.commit()
        database.refresh(new_match)
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('matches table with new match')

    try:
        _ = await populate_predictions(
            database,
            match_id = new_match.match_id,
            match_date = new_match.match_date,
        )
    except:
        # Delete match if predictions could not be created
        try:
            new_match_id = database.query(
                func.max(MatchModel.match_id)
            ).first()[0]
            _ = await delete_match_by_id(new_match_id, database)
        except:
            raise COULD_NOT_UPDATE_EXCEPTION('matches table when deleting match')
        raise COULD_NOT_UPDATE_EXCEPTION('predictions table when populating with new predictions')
    
    return new_match

async def get_match_by_id(match_id: int, database: Session) -> MatchModel:
    '''
    Gets a specified match from the database based on its id
    '''
    match = database.query(MatchModel).get(match_id)
    if not match:
        raise MATCH_NOT_FOUND_EXCEPTION
    return match

async def update_match_by_id(match_id: int, match_score: schema.MatchScore, database: Session) -> MatchModel:
    '''
    Updates a specified match's score
    '''
    match = await get_match_by_id(match_id, database)
    try:
        match.home_goals = match_score.home_goals
        match.away_goals = match_score.away_goals

        database.commit()
    except:
        raise COULD_NOT_UPDATE_EXCEPTION(f'match score for match {match_id}')
    return match

async def get_multiple_matches_by_id(match_ids: List[int], database: Session) -> List[MatchModel]:
    '''
    Retrieves multiple matches from db based on match id
    '''
    matches = database.query(MatchModel).filter(MatchModel.match_id.in_(match_ids)).all()
    if not matches:
        raise MATCH_NOT_FOUND_EXCEPTION
    return matches

async def update_multiple_matches_by_id(match_scores: List[schema.MatchScoreWithId], database: Session) -> List[MatchModel]:
    '''
    Updates multiple matches based on ids
    '''
    home_goals_payload = {ms.match_id:ms.home_goals for ms in match_scores}
    away_goals_payload = {ms.match_id:ms.away_goals for ms in match_scores}

    matches_query = database.query(
        MatchModel
    ).filter(
        MatchModel.match_id.in_(home_goals_payload)
    )

    try:
        matches_query.update({
            MatchModel.home_goals: case(
                home_goals_payload,
                value = MatchModel.match_id
            ),
            MatchModel.away_goals: case(
                away_goals_payload,
                value = MatchModel.match_id
            )
        })
    except:
        raise COULD_NOT_UPDATE_EXCEPTION('matches table with new matches')

    database.commit()

    return matches_query.all()

async def delete_match_by_id(match_id: int, database: Session) -> None:
    database.query(MatchModel).filter(MatchModel.match_id == match_id).delete()
    database.commit()