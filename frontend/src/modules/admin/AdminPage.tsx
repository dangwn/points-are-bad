import React, { useState } from 'react';
import { useRouter } from 'next/router';
import { useQueries } from 'react-query';

import withUser from '../auth/withUser';
import MatchAdmin from './MatchAdmin';
import Loading from '../shared/Loading';
import Error from '../shared/Error';
import { getFullMatches } from '../../lib/adminRequests';
import type { MatchWithId } from '../../types/match';

interface AdminPageProps {
  username: string;
  isAdmin: boolean
}

interface NavigationButtonProps {
  text: string;
  onClick: () => void
}

const NavigationButton: React.FC<NavigationButtonProps> = ({text, onClick}) => {
  return <button onClick={onClick}>{text}</button>
}

const AdminPage: React.FC<AdminPageProps> = ({username, isAdmin}) => {
  const [adminState, setAdminState] = useState<string>('')
  const router = useRouter();

  if (!isAdmin) {
    return <Error />
  }


  return (
    <div>
      <NavigationButton text='PAB Home' onClick={() => router.push('/')} />
      <NavigationButton text='Users' onClick={() => setAdminState('user')} />
      <NavigationButton text='Matches' onClick={() => setAdminState('match')} />
      {
        (adminState === 'match') &&
        <MatchAdmin />
      }
    </div>
  )
};

export default withUser(AdminPage);
