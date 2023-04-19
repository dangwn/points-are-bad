import React, { useState } from 'react';
import { useRouter } from 'next/router';

import withUser from '../auth/withUser';
import MatchAdmin from './MatchAdmin';
import Error from '../shared/Error';

import { updatePoints } from '@/lib/adminRequests';

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
  const [adminState, setAdminState] = useState<string>('');
  const router = useRouter();

  if (!isAdmin) {
    return <Error />
  }


  return (
    <div>
      <NavigationButton text='PAB Home' onClick={() => router.push('/')} />
      <NavigationButton text='Users' onClick={() => setAdminState('user')} />
      <NavigationButton text='Matches' onClick={() => setAdminState('match')} />
      <button
        onClick={async () => {
          try {
            await updatePoints()
          } catch (e) {
            console.log(e);
          }
        }}
      >
        Update Points
      </button>
      {
        (adminState === 'match') &&
        <MatchAdmin />
      }
    </div>
  )
};

export default withUser(AdminPage);
