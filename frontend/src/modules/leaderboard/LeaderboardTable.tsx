import React from 'react';
import { useInfiniteQuery } from 'react-query';

import { getLeaderboard } from '@/lib/requests';
import type { LeaderboardUser } from '@/types/leaderboard';
import styles from '@/styles/LeaderboardPage.module.css'

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
      <h2 className={styles.tableTitle}>Leaderboard</h2>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}></th>
            <th className={styles.th}>Points</th>
            <th className={styles.th}>Correct Scores</th>
            <th className={styles.th}>Largest Error</th>
            <th className={styles.th}>Position</th>
          </tr>
        </thead>
        <tbody>
          {users.map((leadboardUser, index) => (
            <tr key={index}>
              <td className={styles.tableElement}>{leadboardUser?.user.username}</td>
              <td className={styles.tableElement}>{leadboardUser?.points}</td>
              <td className={styles.tableElement}>{leadboardUser?.correct_scores}</td>
              <td className={styles.tableElement}>{leadboardUser?.largest_error}</td>
              <td className={styles.tableElement}>{leadboardUser?.position ? leadboardUser.position : '-'}</td>
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