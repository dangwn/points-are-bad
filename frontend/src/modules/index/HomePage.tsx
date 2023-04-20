import React from 'react';
import { useQueries } from 'react-query';

import PointsTable from './PointsTable';
import UpcomingMatches from './UpcomingMatches';
import Header from '../shared/Header';
import Loading from '../shared/Loading';
import Error from '../shared/Error';

import { createPositionString } from '@/lib/change';
import { getSessionUserPoints, getUpcomingMatches } from '@/lib/requests';
import type { SessionUserPoints, LeaderboardPoints } from '@/types/points';
import type { MatchWithoutGoals } from '@/types/match';

interface HomePageProps {}

const HomePage: React.FC<HomePageProps> = ({}) => {
  const queryResults = useQueries([
    {
      queryKey: 'sessionUserPoints',
      queryFn: getSessionUserPoints,
      refetchOnWindowFocus: false
    },
    {
      queryKey: 'upcomingMatches',
      queryFn: getUpcomingMatches,
      refetchOnWindowFocus: false
    }
  ]);
  
  if (queryResults.some((result) => result.isLoading)) {
    return <Loading />
  };

  if (queryResults.some((result) => result.isError)) {
    return <Error />
  };

  const sessionUserPoints: LeaderboardPoints = queryResults[0].data || {
    points: 0,
    correct_scores: 0,
    largest_error: 0,
    position: null,
    user: {
      username: '',
      is_admin: false
    }
  };
  const upcomingMatches: MatchWithoutGoals[] = queryResults[1].data || [];

  return (
    <div>
      <Header 
        isAdmin={sessionUserPoints.user.is_admin}
      />
      <PointsTable 
        username={`${sessionUserPoints.user.username} | ${createPositionString(sessionUserPoints.position)}`}
        points={sessionUserPoints.points}
        correctScores={sessionUserPoints.correct_scores}
        largestError={sessionUserPoints.largest_error}
      />
      <UpcomingMatches 
        matchData={upcomingMatches}
      />
    </div>
  )
}

export default HomePage;