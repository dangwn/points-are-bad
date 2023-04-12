import { useMutation, useQueryClient } from 'react-query';
import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/header/Header.module.css';

interface HeaderProps {
  username: string,
  isAdmin: boolean
}

const Header: React.FC<HeaderProps> = ({username, isAdmin}) => {
  const [showDropdown, setShowDropdown] = useState<boolean>(false);
  const queryClient = useQueryClient();
  const router = useRouter();

  const handleHome = useMutation(async () => {
    router.push('/')
  })

  const handleLeaderboard = useMutation(async () => {
    router.push('/leaderboard')
  })

  const handlePredictions = useMutation(async () => {
    router.push('/predictions')
  })

  const handleLogout = useMutation(async () => {
    await fetch(`${API_HOST}/auth/login/`, {
      method: 'DELETE',
      credentials: 'include',
    });

    queryClient.removeQueries();

    router.push('/login');
  });

  const handleDeleteUser = useMutation(async () => {
    const accessToken = localStorage.getItem('access_token');
    await fetch(`${API_HOST}/user/`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      },
      credentials: 'include',
    });

    queryClient.removeQueries();

    router.push('/login');
  })

  const handleDropdownClick = () => {
    setShowDropdown(!showDropdown);
  };

  return (
    <header className={styles.header}>
      <div className={styles.title}>
        <Link href="/">
          Points Are Bad
        </Link>
      </div>
      <div className={styles.user}>
        <button className={styles.username} onClick={handleDropdownClick}>
          Menu
        </button>
        {showDropdown && (
          <div className={styles.dropdown}>
            <button className={styles.dropdownButton} onClick={() => handleHome.mutate()}>
              Home
            </button>
            <button className={styles.dropdownButton}>Settings</button>
            <button className={styles.dropdownButton} onClick={() => handlePredictions.mutate()}>
              Predictions
            </button>
            <button className={styles.dropdownButton} onClick={() => handleLeaderboard.mutate()}>
              Leaderboard
            </button>
            { isAdmin ? 
              <button className={styles.dropdownButton}>Admin</button> :
              null
            }
            <button className={styles.dropdownButton} onClick={() => handleLogout.mutate()}>
              Logout
            </button>
            <button className={styles.dropdownButton} onClick={() => handleDeleteUser.mutate()}>
              Delete User
            </button>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;