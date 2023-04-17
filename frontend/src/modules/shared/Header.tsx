import { useMutation, useQueryClient } from 'react-query';
import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';

import { API_HOST } from '../../lib/constants';
import { deleteAccessToken } from '../../lib/accessToken';
import styles from '../../styles/header/Header.module.css';

interface DropdownButtonProps {
  text: string;
  onClick: () => void;
}
interface HeaderProps {
  isAdmin: boolean;
}

const DropdownButton: React.FC<DropdownButtonProps> = ({text, onClick}) => {
  return <button className={styles.dropdownButton} onClick={onClick}>{text}</button> 
}

const Header: React.FC<HeaderProps> = ({ isAdmin }) => {
  const [showDropdown, setShowDropdown] = useState<boolean>(false);
  const queryClient = useQueryClient();
  const router = useRouter();

  const handleLogout = useMutation(async () => {
    await fetch(`${API_HOST}/auth/login/`, {
      method: 'DELETE',
      credentials: 'include',
    });

    queryClient.removeQueries();
    deleteAccessToken();
    router.push('/login');
  });

  const handleDropdownClick = () => {
    setShowDropdown(!showDropdown);
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
            <DropdownButton text='Home' onClick={() => router.push('/')} />
            <DropdownButton text='Settings' onClick={() => router.push('/settings')} />
            <DropdownButton text='Predictions' onClick={() => router.push('/predictions')} />
            <DropdownButton text='Leaderboard' onClick={() => router.push('/leaderboard')} />
            {
              isAdmin && (
              <DropdownButton text='Admin' onClick={() => router.push('/admin')} />
            )}
            <DropdownButton text='Logout' onClick={() => handleLogout.mutate()} />
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;