import { signOut } from "next-auth/react";
import type { Session } from "next-auth";

import LeaderBoard from "../leaderboard/Leaderboard";

import { Inter } from '@next/font/google';

const inter = Inter({ subsets: ['latin'] });

interface userHomePageProps {
  session: Session
}

const UserHomePage: React.FC<userHomePageProps> = ({ session }) => {
  return (
    <>
      <h1 className={inter.className}>{`Welcome ${session.user.name}`}</h1>
      <button onClick={() => signOut({
        callbackUrl: '/'
      })}>Sign Out</button>
      <LeaderBoard leagueId={null}/>
    </>
  )
};
export default UserHomePage;