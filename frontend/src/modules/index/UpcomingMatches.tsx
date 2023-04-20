import React from 'react';

import { MatchWithoutGoals } from '@/types/match';
import styles from '@/styles/HomePage.module.css'

interface UpcomingMatchesProps {
  matchData: MatchWithoutGoals[]
}

const UpcomingMatches: React.FC<UpcomingMatchesProps> = ({ matchData }) => {
  return (
    <div className={styles.upcomingMatchesContainer}>
      <div className={styles.upcomingMatchesTableTitle}>Upcoming Fixtures</div>
      <table className={styles.upcomingMatchesTable}>
        <thead>
          <tr>
            <th className={styles.upcomingMatchesTableElement}>Home</th>
            <th className={styles.upcomingMatchesTableElement}>Match Date</th>
            <th className={styles.upcomingMatchesTableElement}>Away</th>
          </tr>
        </thead>
        <tbody>
          {matchData.map((match, index) => (
            <tr key={index}>
              <td className={styles.upcomingMatchesTableElement}>{match.home}</td>
              <td className={styles.upcomingMatchesTableElement}>{match.match_date}</td>
              <td className={styles.upcomingMatchesTableElement}>{match.away}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default UpcomingMatches;