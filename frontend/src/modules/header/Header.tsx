import React, { useState, useEffect } from "react";
import Image from "next/image";

import HeaderUserStat from "./HeaderUserStat";
import { getUserPoints } from "@/lib/requests/user";

import type { Session } from "next-auth";

import headerStyles from "../../styles/header.module.css";
import { Inter } from '@next/font/google';

const inter = Inter({ subsets: ['latin'] });

interface headerProps {
  session: Session
}

const Header: React.FC<headerProps> = ({ session }) => {
  const [ points, setPoints ] = useState<number>(0);
  const [ correctScores, setCorrectScores ] = useState<number>(0);
  const [ largestError, setLargestError ] = useState<number>(0);
  const [ globalPosition, setGlobalPosition ] = useState<number>(0);
  
  useEffect(() => {
    const setUserPointsData = async (): Promise<void> => {
      const userPointsData = await getUserPoints(
        session.user.accessToken, 
        session.user.provider
      );

      setPoints(userPointsData.points);
      setCorrectScores(userPointsData.correct_scores);
      setLargestError(userPointsData.largest_error);
      setGlobalPosition(userPointsData.position);
    };
    setUserPointsData();
  }, []);

  return (
    <>
      <div className={`${headerStyles.header} ${inter.className}`}>
        <HeaderUserStat statKey={"Points"} value={points}/>
        <HeaderUserStat statKey={"Correct Scores"} value={correctScores}/>
        <HeaderUserStat statKey={"Largest Error"} value={largestError}/>
        <HeaderUserStat statKey={"Global Position"} value={globalPosition}/>
        <h2 className={headerStyles.headerUserName}>{session.user.name}</h2>
        <Image className={headerStyles.headerAvatar} src={session.user.image} alt="User avatar" width={200} height={200}/>
      </div>
    </>
  )
};

export default Header;