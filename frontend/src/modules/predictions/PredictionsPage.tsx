import React from 'react';
import { useQuery } from 'react-query';

import Header from '../shared/Header';
import PredictionsTable from './PredictionsTable';
import ResultsTable from './ResultsTable';
import Loading from '../shared/Loading';
import Error from '../shared/Error';
import withUser from '../auth/withUser';

import { QUERY_OPTIONS } from '@/lib/constants';
import { getUserPredictions } from '@/lib/requests';
import type { UserPrediction } from '@/types/predictions';

interface PredictionsPageProps {
  username: string,
  isAdmin: boolean
}

const PredictionsPage: React.FC<PredictionsPageProps> = ({username, isAdmin}) => {
  const { data, isLoading, isError } = useQuery<UserPrediction[]>(
    'userPredictions',
    getUserPredictions,
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <Loading />;
  };

  if (isError) {
    return <Error />
  };

  const today = new Date();
  const todayDateString = today.toISOString().slice(0, 10);

  const upcomingPredictions: UserPrediction[] = data?.filter(
    (pred) => (pred.match.match_date > todayDateString)
  ) || [];
  const previousPredictions: UserPrediction[] = data?.filter(
    (pred) => (pred.match.match_date <= todayDateString)
  ) || [];

  return (
    <div>
      <Header 
        isAdmin={isAdmin}
      />
      <PredictionsTable predictions={upcomingPredictions}/>
      <ResultsTable predictions={previousPredictions} />
    </div>
  )
};

export default withUser(PredictionsPage);