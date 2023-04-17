import React from 'react';
import { useInfiniteQuery } from 'react-query';

import { getLeaderboard } from '../../lib/requests';
import type { LeaderboardUser } from '../../types/leaderboard';
import styles from '../../styles/leaderboard/LeaderboardTable.module.css'

interface LeaderboardTableProps {
  limit: number
}

const LeaderboardTable: React.FC<LeaderboardTableProps> = ({ limit }) => {
  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useInfiniteQuery(
    'leaderboard',
    ({ pageParam = 0 }) => getLeaderboard(pageParam, limit),
    {
      getNextPageParam: lastPage => {
        const hasMore = lastPage.hasMore;
        
        return hasMore ? (lastPage.pageIndex + 1) : undefined;
      },
      refetchOnWindowFocus: false
    }
  );

  const users: LeaderboardUser[] = data ? data.pages.flatMap(page => page.data) : [];
  return (
    <div className={styles.container}>
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Username</th>
            <th>Points</th>
            <th>Correct Scores</th>
            <th>Largest Error</th>
            <th>Position</th>
          </tr>
        </thead>
        <tbody>
          {users.map((leadboardUser, index) => (
            <tr key={index}>
              <td>{leadboardUser?.user.username}</td>
              <td>{leadboardUser?.points}</td>
              <td>{leadboardUser?.correct_scores}</td>
              <td>{leadboardUser?.largest_error}</td>
              <td>{leadboardUser?.position ? leadboardUser.position : '-'}</td>
            </tr>
          ))}
        </tbody>
      </table>
      {hasNextPage && (
        <button onClick={() => fetchNextPage()} disabled={isFetchingNextPage} className={styles.button}>
          {isFetchingNextPage ? 'Loading more...' : 'Load More'}
        </button>
      )}
    </div>
  );
}

export default LeaderboardTable;