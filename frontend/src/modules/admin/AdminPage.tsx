import React, { useState } from 'react';
import { useRouter } from 'next/router';

import MatchAdmin from './MatchAdmin';
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
    return <>404</>
  }

  const Matches: MatchWithId[] = [
    {
      match_id: 1,
      match_date: '2023-01-01',
      home: 'England',
      away: 'Germany',
      home_goals: 2,
      away_goals: 0
    },
    {
      match_id: 2,
      match_date: '2023-01-01',
      home: 'Italy',
      away: 'France',
      home_goals: null,
      away_goals: null
    }
  ]

  return (
    <div>
      <NavigationButton text='PAB Home' onClick={() => router.push('/')} />
      <NavigationButton text='Users' onClick={() => setAdminState('user')} />
      <NavigationButton text='Matches' onClick={() => setAdminState('match')} />
      {
        (adminState === 'match') &&
        <MatchAdmin inputMatchData={Matches} />
      }
    </div>
  )
};

export default AdminPage;
