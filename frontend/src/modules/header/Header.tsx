import { useMutation, useQueryClient, useQuery } from 'react-query';
import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';

import { API_HOST, QUERY_OPTIONS } from '../../lib/constants';
import styles from '../../styles/header/Header.module.css';

type SessionUser = {
  username: string,
  is_admin: boolean
}

const Header: React.FC = () => {
  const [showDropdown, setShowDropdown] = useState<boolean>(false);
  const queryClient = useQueryClient();
  const router = useRouter();

  const { data, isLoading, isError } = useQuery<SessionUser>(
    'user',
    async () => {
      const authToken = localStorage.getItem('access_token');
  
      if (authToken === 'undefined') {
        throw new Error('Auth token could not be found');
      };
      const response = await fetch(`${API_HOST}/user/`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
    
      if (!response.ok){
        throw new Error('Error fetching user data');
      };
      return response.json();
    },
    QUERY_OPTIONS
  );

  const handleLogout = useMutation(async () => {
    await fetch(`${API_HOST}/auth/login/`, {
      method: 'DELETE',
      credentials: 'include',
    });

    queryClient.removeQueries();

    router.push('/login');
  });

  const handleDropdownClick = () => {
    setShowDropdown(!showDropdown);
  };

  const isAdmin: boolean = data?.is_admin || false

  if (isLoading) {
    return <div>Loading...</div>;
  };

  if (isError) {
    return <div>Error fetching home page data</div>
  };

  return (
    <header className={styles.header}>
      <div className={styles.title}>
        <Link className={styles.a} href="/">
          Points Are Bad
        </Link>
      </div>
      <div className={styles.user}>
        <button className={styles.username} onClick={handleDropdownClick}>
          Menu
        </button>
        {showDropdown && (
          <div className={styles.dropdown}>
            <button 
              className={styles.dropdownButton} 
              onClick={() => router.push('/')}>
              Home
            </button>
            <button 
              className={styles.dropdownButton}
              onClick={() => router.push('/settings')}
            >
              Settings
            </button>
            <button 
              className={styles.dropdownButton} 
              onClick={() => router.push('/predictions')}>
              Predictions
            </button>
            <button 
              className={styles.dropdownButton} 
              onClick={() => router.push('/leaderboard')}
            >
              Leaderboard
            </button>
            { isAdmin ? 
              <button className={styles.dropdownButton}>Admin</button> :
              null
            }
            <button className={styles.dropdownButton} onClick={() => handleLogout.mutate()}>
              Logout
            </button>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;