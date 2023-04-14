import React from 'react';

import Header from '../modules/header/Header';
import DeleteUser from '../modules/settings/DeleteUser';
import withAuth from '../modules/auth/withAuth';

const Leaderboard: React.FC = () => {
  return (
    <div>
      <Header />
      <DeleteUser />
    </div>
  )
}

export default withAuth(Leaderboard);