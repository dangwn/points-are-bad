import React from 'react'
import { useSession, signOut, signIn } from "next-auth/react";
import { useRouter } from 'next/router';

import SignMeUpForm from '@/modules/index/SignMeUpForm';
import LandingPage from '@/modules/index/LandingPage';
import UserHomePage from '@/modules/index/UserHomePage';


interface homeProps {};

const Home: React.FC<homeProps> = ({}) => {
  const { data: session, status } = useSession();

  if (session?.user.newUser) {
    return <SignMeUpForm session={session}/>
  };
  if (session) {
    return <UserHomePage session={session}/>
  }
  if (status === 'loading'){
    return <h1>Loading...</h1>
  }
  
  return <LandingPage />
};

export default Home;