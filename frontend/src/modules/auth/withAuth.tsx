import Cookies from 'js-cookie';
import { useRouter } from 'next/router';
import React, { useEffect } from 'react';

import { setAccessToken } from '../../lib/accessToken';
import { refreshAccessToken } from '../../lib/requests';
import type { Token } from '../../types/token';

interface WithAuthProps {}

const withAuth = (WrappedComponent: React.FC<WithAuthProps>) => {
  const AuthWrapper: React.FC = (props) => {
    const router = useRouter();

    const csrfToken: string|undefined = Cookies.get('X-CSRF-Token');
    const refreshToken: string|undefined = Cookies.get('X-Refresh-Token');
    const deleteSession = (): void => {
      try {
        Cookies.remove('X-CSRF-Token');
        Cookies.remove('X-Refresh-Token');
      } catch {};
      router.push('/login');
    };

    useEffect(() => {
      if (csrfToken === undefined || refreshToken === undefined) {
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