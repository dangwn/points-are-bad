import React from 'react';

import withAuth from '../modules/auth/withAuth';
import Header from '../modules/header/Header';
import PointsTable from '../modules/index/PointsTable';
import UpcomingMatches from '../modules/index/UpcomingMatches';

const HomePage: React.FC = () => {  
  return (
    <div>
      <Header />
      <PointsTable />
      <UpcomingMatches />
    </div>
  );
};

export default withAuth(HomePage);