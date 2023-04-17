import React from 'react';
import { useQueries } from 'react-query';

import PointsTable from './PointsTable';
import UpcomingMatches from './UpcomingMatches';
import Header from '../shared/Header';
import Loading from '../shared/Loading';

import { getSessionUserPoints, getUpcomingMatches } from '../../lib/requests';
import { SessionUserPoints } from '../../types/points';
import { MatchWithoutGoals } from '../../types/match';

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
    return <div>Error fetching data</div>
  };

  const sessionUserPoints: SessionUserPoints = queryResults[0].data;
  const upcomingMatches: MatchWithoutGoals[] = queryResults[1].data || [];

  return (
    <div>
      <Header 
        isAdmin={sessionUserPoints.user.is_admin}
      />
      <PointsTable 
        username={sessionUserPoints.user.username}
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