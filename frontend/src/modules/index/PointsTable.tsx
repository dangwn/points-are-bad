import React from 'react';
import { useQuery } from 'react-query';

import styles from '../../styles/index/IndexTables.module.css';
import { API_HOST, QUERY_OPTIONS } from '@/lib/constants';

type SessionUserPoints = {
  points: number,
  correct_scores: number,
  largest_error: number,
  user: {
    username: string,
    is_admin: boolean
  }
}

const PointsTable: React.FC = () => {
  const { data, isLoading, isError } = useQuery<SessionUserPoints>(
    'points',
    async () => {
      const authToken = localStorage.getItem('access_token');
  
      if (authToken === 'undefined') {
        throw new Error('Auth token could not be found');
      };
      const response = await fetch(`${API_HOST}/points/`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });
    
      if (!response.ok){
        throw new Error('Error fetching points data');
      };
      return response.json();
    },
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <div>Loading...</div>;
  };

  if (isError) {
    return <div>Error fetching home page data</div>
  };

  const username: string = data?.user.username || ''
  const points: number = (data ? data.points : NaN)
  const correctScores: number = (data ? data.correct_scores : NaN)
  const largestError: number = (data ? data.largest_error : NaN)

  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>{username}</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}>Points</th>
            <th className={styles.th}>Correct Scores</th>
            <th className={styles.th}>Largest Error</th>
          </tr>
        </thead>
        <tbody>
          <tr className={styles.tr}>
            <td className={styles.td}>{points}</td>
            <td className={styles.td}>{correctScores}</td>
            <td className={styles.td}>{largestError}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default PointsTable;
