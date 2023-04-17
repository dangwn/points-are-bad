import { useRouter } from 'next/router';
import React, { useEffect } from 'react';
import Cookies from 'js-cookie';

const withLogin = (WrappedComponent: React.FC) => {
  const LoginWrapper: React.FC = (props) => {
    const router = useRouter();

    useEffect(() => {
      const csrfToken = Cookies.get('X-CSRF-Token');
      const refreshToken = Cookies.get('X-Refresh-Token');

      if (csrfToken !== undefined && refreshToken !== undefined) {
        router.push('/');
      };

    }, []);

    return <WrappedComponent {...props} />
  };

  return LoginWrapper;
};

export default withLogin;