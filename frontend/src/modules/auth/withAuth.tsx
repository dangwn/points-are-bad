import Cookies from 'js-cookie';
import { useRouter } from 'next/router';
import React, { useEffect } from 'react';

import { setAccessToken } from '../../lib/accessToken';
import { refreshAccessToken } from '../../lib/requests';
import type { Token } from '../../types/token';

import jwt_decode from 'jwt-decode';
import { getAccessToken } from '../../lib/accessToken';

interface WithAuthProps {}

type DecodedAccessToken = {
  iat: number,
  exp: number,
  is_admin: boolean,
  username: string,
  sub: string
}

/*
 * Component wrapper to keep access token updated
 * Wrapper checks if the expire time of the token is after the current time
 *  - If it is, just return component
 *  - If not, hit refresh endpoint to get a new refresh token
 *      If an error is thrown, log the user out and send them to login
 */
const withAuth = (WrappedComponent: React.FC<WithAuthProps>) => {
  const AuthWrapper: React.FC = (props) => {
    /*
     * Check if token is valid
     * The if statement is here to stop the server trying to access
     *    localStorage in getAccessToken while not on client side
     */
    if (typeof window !== 'undefined') {
      try {
        const accessToken: string = getAccessToken();
        const decodedToken: DecodedAccessToken = jwt_decode(accessToken);
        if (decodedToken['exp'] > Date.now()/1000) {
          return <WrappedComponent {...props}/>
        }
      } catch {}
    }

    const router = useRouter();

    const refreshToken: string|undefined = Cookies.get('X-Refresh-Token');
    const deleteSession = (): void => {
      try {
        Cookies.remove('X-Refresh-Token');
      } catch {};
      router.push('/login');
    };

    // Fetch refresh token, but if it's undefined or the access token could not be
    //    fetched, delete the session
    useEffect(() => {
      if (refreshToken === undefined) {
        deleteSession();
      };

      const injectUserSession = async (): Promise<void> => {
        try {
          const { access_token }: Token = await refreshAccessToken();
          setAccessToken(access_token);
        } catch {
          deleteSession();
        };
      };

      injectUserSession();
    }, []);

    return <WrappedComponent {...props}/>
  };

  return AuthWrapper;
};

export default withAuth;