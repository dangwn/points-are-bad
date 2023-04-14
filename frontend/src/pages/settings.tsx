import React from 'react';

import Header from '../modules/header/Header';
import DangerZone from '../modules/settings/DangerZone';
import withAuth from '../modules/auth/withAuth';

const Leaderboard: React.FC = () => {
  return (
    <div>
      <Header />
      <DangerZone />
    </div>
  )
}

export default withAuth(Leaderboard);