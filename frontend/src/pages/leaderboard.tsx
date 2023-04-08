import React from "react";
import { useQuery } from "react-query";

import Header from "../modules/header/Header";
import LeaderboardTable from "../modules/leaderboard/LeaderboardTable";
import withAuth from "../modules/auth/withAuth";

import { API_HOST } from "../lib/constants";

type SessionUser = {
  username: string,
  is_admin: boolean
}

const Leaderboard: React.FC = () => {
  const queryOptions = {
    refetchOnWindowFocus: false
  };

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
        throw new Error('Error fetching user data');
      };
      return response.json();
    },
    queryOptions
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
        username={data?.username} 
        isAdmin={data?.is_admin}
      />
      <LeaderboardTable limit={20} />
    </div>
  )
}

export default withAuth(Leaderboard);