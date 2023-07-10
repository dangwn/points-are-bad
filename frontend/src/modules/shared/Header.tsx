import { logUserOut } from '@/lib/requests';
import Link from 'next/link';
import { useRouter } from 'next/router';
import React, { useState } from 'react';
import { useMutation, useQueryClient } from 'react-query';

import { deleteAccessToken } from '@/lib/accessToken';
import { CLIENT_HOST } from '@/lib/constants';
import styles from '@/styles/Header.module.css';

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
    await logUserOut();

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
        <button className={styles.menu} onClick={handleDropdownClick}>
          Menu
        </button>
        {showDropdown && (
          <div className={styles.dropdown}>
            <button className={styles.dropdownButtonTop} onClick={() => router.push('/')}>Home</button> 
            <DropdownButton text='Settings' onClick={() => router.push('/settings')} />
            <DropdownButton text='Predictions' onClick={() => router.push('/predictions')} />
            <DropdownButton text='Leaderboard' onClick={() => router.push('/leaderboard')} />
            {
              isAdmin && (
              <DropdownButton text='Admin' onClick={() => window.open(`${CLIENT_HOST}/admin`,'_blank')} />
            )}
            <button className={styles.dropdownButtonBottom} onClick={() => handleLogout.mutate()}>Logout</button> 
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;