import React from "react";

import leaderBoardStyles from "../../styles/leaderBoard.module.css";

import type UserPointsData from "@/types/UserPointsData";

interface leaderBoardRowProps {
  data: UserPointsData
};

const LeaderBoardRow: React.FC<leaderBoardRowProps> = ({data}) => {
  return (
    <div className={leaderBoardStyles.leaderBoardRow}>
      <div>{data.display_name}</div>
    </div>
  )
};

export default LeaderBoardRow;