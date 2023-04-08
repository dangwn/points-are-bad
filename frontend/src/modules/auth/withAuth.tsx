import { useRouter } from 'next/router';
import React, { useEffect } from 'react';
import Cookies from 'js-cookie';

import { API_HOST } from '../../lib/constants';

type Token = {
  access_token: string,
  token_type: string
}

const withAuth = (WrappedComponent: React.FC) => {
  const AuthWrapper: React.FC = (props) => {
    const router = useRouter();

    useEffect(() => {
      const csrfToken = Cookies.get('X-CSRF-Token');
      const refreshToken = Cookies.get('X-Refresh-Token');

      if (csrfToken === undefined || refreshToken === undefined) {
        try {
          Cookies.remove('X-CSRF-Token');
          Cookies.remove('X-Refresh-Token');
        } catch {};
        router.push('/login');
      }

      const setAccessToken = async () => {
        const response = await fetch(
          `${API_HOST}/auth/refresh/`,
          {
            method: 'POST',
            headers: {
              'accept': 'application/json'
            },
            credentials: 'include'
          }
        );
        if (!response.ok) {
          router.push('/login');
        };
        const data: Token = await response.json();
        localStorage.setItem('access_token', data.access_token);
      };

      setAccessToken();

    }, []);

    return <WrappedComponent {...props} />
  };

  return AuthWrapper;
};

export default withAuth;