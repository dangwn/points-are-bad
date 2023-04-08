import React from 'react';
import { useQuery } from 'react-query';

import withAuth from '../modules/auth/withAuth';
import Header from '../modules/header/Header';
import PointsTable from '../modules/index/PointsTable';
import UpcomingMatches from '../modules/index/UpcomingMatches';

import { API_HOST, QUERY_OPTIONS } from '../lib/constants';

type SessionUser = {
  username: string,
  is_admin: boolean
}

type SessionUserPoints = {
  points: number,
  correct_scores: number,
  largest_error: number,
  user: SessionUser
}

const HomePage: React.FC = () => {
  const { data, isLoading, isError } = useQuery<SessionUserPoints>(
    'points',
    async () => {
      const authToken = localStorage.getItem('access_token');
  
      if (authToken === 'undefined') {
        throw new Error('Auth token could not be found');
      };
      const response = await fetch(`${API_HOST}/points/`, {
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
    return <div>Error fetching home page data</div>
  };

  return (
    <div>
      <Header 
        username={data?.user.username} 
        isAdmin={data?.user.is_admin}
      />
      <PointsTable 
        username={data?.user.username} 
        points={data?.points} 
        correctScores={data?.correct_scores} 
        largestError={data?.largest_error} 
      />
      <UpcomingMatches />
    </div>
  );
};

export default withAuth(HomePage);