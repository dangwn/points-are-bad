import React from 'react';
import { useInfiniteQuery } from 'react-query';

import { API_HOST } from '../../lib/constants';

import styles from '../../styles/leaderboard/LeaderboardTable.module.css'

interface LeaderboardTableProps {
  limit: number
}

type LeaderboardUser = {
  user: {
    username: string
  },
  points: number,
  correct_scores: number,
  largest_error: number,
  position: number|null
}

type LeaderboardApiResponse = {
  data: LeaderboardUser[];
  hasMore: boolean;
  pageIndex: number
}

const fetchLeaderboard = async (pageIndex: number, pageSize: number): Promise<LeaderboardApiResponse> => {
  const offset = pageIndex * pageSize;
  const response = await fetch(
    `${API_HOST}/points/leaderboard?limit=${pageSize}&offset=${offset}`
  )
  
  if (!response.ok) {
    return {data: [], hasMore: false, pageIndex: pageIndex}
  };
  const data = await response.json();
  if (data.length < pageSize) {
    return {data: data, hasMore: false, pageIndex: pageIndex};
  }
  return {data: data, hasMore: true, pageIndex: pageIndex}
}

const LeaderboardTable: React.FC<LeaderboardTableProps> = ({ limit }) => {

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage } = useInfiniteQuery(
    'leaderboard',
    ({ pageParam = 0 }) => fetchLeaderboard(pageParam, limit),
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