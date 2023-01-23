import NextAuth, { NextAuthOptions } from "next-auth";
import GoogleProvider from "next-auth/providers/google";
import GithubProvider from "next-auth/providers/github";
import { 
  GITHUB_CLIENT_ID, 
  GITHUB_CLIENT_SECRET, 
  GOOGLE_CLIENT_ID, 
  GOOGLE_CLIENT_SECRET 
} from "@/lib/constants";
import { startUserSession } from "@/lib/requests/auth";

import type { Account } from "next-auth";

const getAccountAccessToken = (account: Account): string => {
  let accessToken = '';

  switch (account.provider) {
    case 'google':
      //@ts-ignore
      accessToken = account.id_token;
      break;
    case 'github':
      accessToken = account.access_token || '';
      break;
    default:
      accessToken = '';
  };
  
  return accessToken;
};

export const authOptions: NextAuthOptions = {
  providers: [
    GithubProvider({
      clientId: GITHUB_CLIENT_ID,
      clientSecret: GITHUB_CLIENT_SECRET,
    }),
    GoogleProvider({
      clientId: GOOGLE_CLIENT_ID,
      clientSecret: GOOGLE_CLIENT_SECRET
    })
  ],
  callbacks: {
    signIn: async({}) => {
      return true;
    },
    jwt: async ({ token, account }) => {
      if (!account) {
        return token;
      };

      // Add provider to token
      token.provider = account.provider;

      token.providerAccessToken = getAccountAccessToken(account);
      
      // Try to start a user session
      // If the user is new, tag the session to redirect to signup
      try {
        //@ts-ignore
        const userData = await startUserSession(token.providerAccessToken, token.provider);
        if (userData.displayName === ''){
          token.newUser = true;
        } else {
          token.newUser = false;
        };
      } catch {
        token.newUser = true;
      };
      
      return token;
    },
    session: async ({ session, token }) => {
      // Add account provider, access token, and new user tag to session
      session.user.accessToken = token.providerAccessToken;
      session.user.provider = token.provider;
      session.user.newUser = token.newUser;
      
      return session;
    },
  },
};

export default NextAuth(authOptions);