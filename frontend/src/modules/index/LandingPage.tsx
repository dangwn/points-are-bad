import React from "react";
import { useRouter } from "next/router";
import { Inter } from '@next/font/google';

const inter = Inter({ subsets: ['latin'] });

const LandingPage: React.FC = () => {
  const { push } = useRouter();
  
  return ( 
    <>
      <h1 className={inter.className}> 
        Welcome to Points Are Bad
      </h1>
      <button onClick={() => {
        push({
          pathname: '/login'
        })
      }}>Sign in</button>

    </>
  )
}

export default LandingPage;