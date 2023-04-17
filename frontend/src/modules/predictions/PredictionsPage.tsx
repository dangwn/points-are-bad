import React from 'react';
import { useQuery } from 'react-query';

import Header from '../shared/Header';
import PredictionsTable from './PredictionsTable';
import Loading from '../shared/Loading';
import Error from '../shared/Error';
import withUser from '../auth/withUser';

import { QUERY_OPTIONS } from '../../lib/constants';
import { getUpcomingUserPredictions } from '../../lib/requests';
import type { UserPrediction } from '../../types/predictions';

interface PredictionsPageProps {
  username: string,
  isAdmin: boolean
}

const PredictionsPage: React.FC<PredictionsPageProps> = ({username, isAdmin}) => {
  const { data, isLoading, isError } = useQuery<UserPrediction[]>(
    'userPredictions',
    getUpcomingUserPredictions,
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <Loading />;
  };

  if (isError) {
    return <Error />
  };

  const userPredictions: UserPrediction[] = data || []

  return (
    <div>
      <Header 
        isAdmin={isAdmin}
      />
      <PredictionsTable predictions={userPredictions}/>
    </div>
  )
};

export default withUser(PredictionsPage);