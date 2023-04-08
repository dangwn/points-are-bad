import { useQuery } from 'react-query';
import React from 'react';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/index/UnderConstruction.module.css'

type SessionUser = {
  username: string,
  is_admin: boolean
}

const getSessionUser = async (): Promise<SessionUser> => {
  const authToken = localStorage.getItem('access_token');
  
  if (authToken === undefined) {
    throw new Error('Oops')
  }
  const response = await fetch(`${API_HOST}/user/`, {
    headers: {
      Authorization: `Bearer ${authToken}`,
    },
  });

  if (!response.ok){
    throw new Error('Could not retrieve session user')
  };

  return response.json();
}

const UserButton: React.FC = () => {
  const { isLoading, data } = useQuery<SessionUser>(
    'user', 
    () => getSessionUser()
  );

  return (
    <div className={styles.container}>
      <button className={styles.button} disabled={isLoading}>
        {
          isLoading ?
          'Loading...' : 
          `Hello ${data?.username}. You are ${data?.is_admin ? "" : "not"} an admin`
        }
      </button>
    </div>
  );
}

export default UserButton;
