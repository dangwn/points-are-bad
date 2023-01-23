import { signOut } from "next-auth/react";
import { Inter } from '@next/font/google';

import type { Session } from "next-auth";

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
    </>
  )
};
export default UserHomePage;