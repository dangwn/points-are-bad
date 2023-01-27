import React from 'react'
import { useSession } from "next-auth/react";

import SignMeUpForm from '@/modules/index/SignMeUpForm';
import LandingPage from '@/modules/index/LandingPage';
import UserHomePage from '@/modules/index/UserHomePage';
import Header from '@/modules/header/Header';

interface homeProps {};

const Home: React.FC<homeProps> = ({}) => {
  const { data: session, status } = useSession();

  if (session?.user.newUser) {
    return <SignMeUpForm session={session}/>
  };
  if (session) {
    return (
      <> 
        <Header session={session} />
        <UserHomePage session={session}/>
      </>
    )
  }
  if (status === 'loading'){
    return <h1>Loading...</h1>
  }
  
  return <LandingPage />
};

export default Home;