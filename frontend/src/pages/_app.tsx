import React from 'react';
import type { AppProps } from 'next/app';
import { SessionProvider } from "next-auth/react";

import type { Session } from "next-auth";

import "../styles/loginPage.css";

const App: React.FC<AppProps> = ({ 
  Component, pageProps: { session, ...pageProps }
}: AppProps<{ session: Session}>) => {
  return (
    <SessionProvider session={session}>
      <Component {...pageProps} />
    </SessionProvider>
  )
};

export default App;