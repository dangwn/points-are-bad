from typing import List, Tuple
from sqlalchemy.orm import Session
from sqlalchemy import func, and_

from fastapi import HTTPException
from points.models import PlayerPoints as PlayerPointsModel
from user.models import User as UserModel
from points.schema import DayScore, DayScoreWithUserId, PlayerPosition, PlayerPositionWithUsername
from http_exceptions import USER_NOT_FOUND_EXCEPION, COULD_NOT_UPDATE_EXCEPTION

PlayerPosition_keys = PlayerPosition.__fields__.keys()
PlayerPositionWithUsername_keys = PlayerPositionWithUsername.__fields__.keys()

async def get_all_player_scores(database: Session) -> List[PlayerPointsModel]:
    '''
    Runs db query to get all users' scores
    '''
    player_scores = database.query(PlayerPointsModel).all()
    return player_scores

async def get_player_score_by_id(user_id: int, database: Session) -> PlayerPointsModel:
    '''
    Runs db query to get particular user's score
    '''
    player_score = database.query(PlayerPointsModel).filter(PlayerPointsModel.user_id == user_id).first()
    if player_score is None:
        raise USER_NOT_FOUND_EXCEPION
    return player_score

async def get_user_position(user_id: int, database: Session) -> PlayerPointsModel:
    '''
    Gets the rank of the given user in the points table
    Ordered by points, descending correct score, largest error
    '''
    subq = database.query(
        PlayerPointsModel, 
        func.rank().over(
            order_by = [
                PlayerPointsModel.points,
                PlayerPointsModel.correct_scores.desc(),
                PlayerPointsModel.largest_error
            ]
        ).label('position')
    ).subquery()

    # Result comes back as a tuple
    query_result = database.query(subq).filter(subq.c.user_id == user_id).first()

    player_position = PlayerPosition(**dict(zip(PlayerPosition_keys, query_result)))

    return player_position

async def update_player_score_by_id(
    user_id: int,
    score: DayScore,
    database: Session
) -> PlayerPointsModel:
    '''
    Updates the given user's points, correct scores and largest error
    '''
    player_score = database.query(PlayerPointsModel).filter(PlayerPointsModel.user_id == user_id).first()
    if player_score is None:
        raise USER_NOT_FOUND_EXCEPION
    
    # Update player score
    player_score.points += score.points
    player_score.correct_scores += score.correct_scores
    if score.largest_error > player_score.largest_error:
        player_score.largest_error = score.largest_error

    # Commit changes to DB
    database.commit()
    
    return player_score

async def update_multiple_scores(
    scores: List[DayScoreWithUserId],
    database: Session
) -> List[PlayerPointsModel]:
    '''
    Updates multiple scores in the database
    '''
    new_scores = []
    for score in scores:
        try:
            player_score = database.query(PlayerPointsModel).filter(PlayerPointsModel.user_id == score.user_id).first()
            if player_score is None:
                raise USER_NOT_FOUND_EXCEPION
            
            # Update player score
            player_score.points += score.points
            player_score.correct_scores += score.correct_scores
            if score.largest_error > player_score.largest_error:
                player_score.largest_error = score.largest_error

            new_scores.append(player_score)
        except HTTPException:
            continue
        except:
            raise COULD_NOT_UPDATE_EXCEPTION(f'points table due to non-http error')
        
    if new_scores != []:
        database.commit()
    return new_scores

async def get_leaderboard(
    table_start: int,
    table_end: int,
    database: Session
):
    '''
    Returns the leaderboard from table start (inc) to table end (exc)
    '''
    subq = database.query(
        PlayerPointsModel, 
        func.rank().over(
            order_by = [
                PlayerPointsModel.points,
                PlayerPointsModel.correct_scores.desc(),
                PlayerPointsModel.largest_error
            ]
        ).label('position')
    ).subquery()

    player_position_result = database.query(
        subq.c.points,
        subq.c.correct_scores,
        subq.c.largest_error,
        subq.c.position,
        UserModel.username
    ).filter(and_(
        subq.c.position >= table_start,
        subq.c.position < table_end,
    )).join(UserModel, subq.c.user_id == UserModel.id, isouter = True).all()

    player_positions = [
        PlayerPositionWithUsername(**dict(zip(PlayerPositionWithUsername_keys, ppr))) for ppr in player_position_result
    ]

    return player_positions
