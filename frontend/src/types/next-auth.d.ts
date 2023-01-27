import NextAuth, { DefaultSession } from "next-auth/next";
import { JWT } from "next-auth/jwt";

declare module 'next-auth' {
  /*
   * Returned by `useSession`, `getSession` and received as a prop on the `SessionProvider` React Context
   */
  interface Session {
    user: {
      provider: string,
      accessToken: string,
      newUser: boolean
    } & DefaultSession['user'],
  };
};

declare module 'next-auth/jwt' {
  interface JWT {
    JWT,
    providerAccessToken: string,
    provider: string,
    newUser: boolean
  }
}