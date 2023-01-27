import React, { useState, useEffect} from "react";
import LeaderBoardRow from "./LeaderboardRow";

import { getGlobalLeaderBoard } from "@/lib/requests/leaderBoard";

import type UserPointsData from "@/types/UserPointsData";

interface leaderBoardProps {
  leagueId: number | null;
};

const LeaderBoard: React.FC<leaderBoardProps> = ({ leagueId }) => {
  // For the time being skipping the league logic
  const [ data, setData ] = useState<UserPointsData[]>([]);

  useEffect(() => {
    const displayData = async () => {
      const leaderBoardData: UserPointsData[] = await getGlobalLeaderBoard(1,10);
      setData(leaderBoardData);
    };
    displayData();
  }, [])


  return (
    <>
      <div>{JSON.stringify(data)}</div>
      <LeaderBoardRow data={data[0]}/>
    </>
  )
};

export default LeaderBoard;