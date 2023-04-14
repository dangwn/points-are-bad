import React from "react";

import Header from "../modules/header/Header";
import LeaderboardTable from "../modules/leaderboard/LeaderboardTable";
import withAuth from "../modules/auth/withAuth";

const Leaderboard: React.FC = () => {
  return (
    <div>
      <Header />
      <LeaderboardTable limit={20} />
    </div>
  )
}

export default withAuth(Leaderboard);