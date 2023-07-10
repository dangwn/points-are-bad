import jwt_decode from 'jwt-decode';
import React from 'react';

import { getAccessToken } from '@/lib/accessToken';

interface WithUserProps {
  username: string;
  isAdmin: boolean
}

type DecodedAccessToken = {
  iat: number,
  exp: number,
  is_admin: boolean,
  username: string,
  sub: string
}

/*
 * Component wrapper to inject user's username and admin status into component
 * Decode the access token to get the username and admin status
 */
const withUser = (WrappedComponent: React.FC<WithUserProps>) => {
  const UserWrapper: React.FC = (props) => {
      if (typeof window !== 'undefined') {
        const accessToken: string = getAccessToken();
        const decodedToken: DecodedAccessToken = jwt_decode(accessToken);
        return <WrappedComponent username={decodedToken['username']} isAdmin={decodedToken['is_admin']} {...props} />
      }
      return <WrappedComponent username={''} isAdmin={false} {...props} />
  };

  return UserWrapper;
}

export default withUser;