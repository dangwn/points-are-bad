import React, { useState } from 'react'
import { useRouter } from 'next/router';
import { useQuery } from 'react-query';
import Cookies from 'js-cookie';

import Loading from '../shared/Loading';
import { QUERY_OPTIONS } from '../../lib/constants';
import { getSessionUser } from '../../lib/requests';
import type { SessionUser } from '../../types/user';

interface WithUserProps {
  username: string;
  isAdmin: boolean
}

const withUser = (WrappedComponent: React.FC<WithUserProps>) => {
  const UserWrapper: React.FC = (props) => {
    const router = useRouter();
    const { data, isLoading, isError } = useQuery<SessionUser>(
      'user',
      getSessionUser,
      QUERY_OPTIONS
    );

    const deleteSession = (): void => {
      try {
        Cookies.remove('X-CSRF-Token');
        Cookies.remove('X-Refresh-Token');
      } catch {};
      router.push('/login');
    };

    const username: string = data?.username || ''
    const isAdmin: boolean = data?.is_admin || false;

    if (isLoading) {
      return <Loading />;
    };
  
    if (isError) {
      deleteSession();
    };

    return <WrappedComponent username={username} isAdmin={isAdmin} {...props} />
  };

  return UserWrapper;
}

export default withUser;