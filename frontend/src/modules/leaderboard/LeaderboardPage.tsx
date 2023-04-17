import React from 'react';

import withUser from '../auth/withUser';
import Header from '../shared/Header';
import LeaderboardTable from './LeaderboardTable';

interface LeaderboardPageProps {
  username: string,
  isAdmin: boolean
};

const LeaderboardPage: React.FC<LeaderboardPageProps> = ({username, isAdmin}) => {

  return (
    <>
      <Header isAdmin={isAdmin}/>
      <LeaderboardTable limit={10}/>
    </>
  )
}

export default withUser(LeaderboardPage);