import Cookies from 'js-cookie';
import { useRouter } from 'next/router';
import React, { useEffect } from 'react';

/*
 * Component wrapper to verify whether a user is logged in or not
 * If the refresh token is present, push them to the home page 
 */
const withLogin = (WrappedComponent: React.FC) => {
  const LoginWrapper: React.FC = (props) => {
    const router = useRouter();

    useEffect(() => {
      const refreshToken = Cookies.get('X-Refresh-Token');

      if (refreshToken !== undefined) {
        router.push('/');
      };

    }, []);

    return <WrappedComponent {...props} />
  };

  return LoginWrapper;
};

export default withLogin;