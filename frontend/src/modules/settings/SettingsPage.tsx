import React from 'react';

import Header from '../shared/Header';
import DangerZone from '../settings/DangerZone';
import withUser from '../auth/withUser';

interface LeaderboardProps {
  username: string;
  isAdmin: boolean
}

const LeaderboardPage: React.FC<LeaderboardProps> = ({username, isAdmin}) => {
  return (
    <div>
      <Header isAdmin={isAdmin}/>
      <DangerZone />
    </div>
  )
}

export default withUser(LeaderboardPage);