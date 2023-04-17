import React from 'react';
import { useQuery } from 'react-query';

import withAuth from '../modules/auth/withAuth';
import Header from '../modules/shared/Header';
import PredictionsTable from '../modules/predictions/PredictionsTable';

import { API_HOST, QUERY_OPTIONS } from '../lib/constants';

type SessionUser = {
  username: string,
  is_admin: boolean
}

const PredictionsPage: React.FC = () => {
  const { data, isLoading, isError } = useQuery<SessionUser>(
    'user',
    async () => {
      const authToken = localStorage.getItem('access_token');
  
      if (authToken === 'undefined') {
        throw new Error('Auth token could not be found');
      };
      const response = await fetch(`${API_HOST}/user/`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
    
      if (!response.ok){
        throw new Error('Error fetching points data');
      };
      return response.json();
    },
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <div>Loading...</div>;
  };

  if (isError) {
    return <div>Error fetching user data</div>
  };

  return (
    <div>
      <Header 
        isAdmin={data?.is_admin}
      />
      <PredictionsTable />
    </div>
  )
};

export default withAuth(PredictionsPage)