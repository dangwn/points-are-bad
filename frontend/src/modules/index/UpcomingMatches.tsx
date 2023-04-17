import React from "react";

import { MatchWithoutGoals } from "../../types/match";
import styles from "../../styles/index/IndexTables.module.css";

interface UpcomingMatchesProps {
  matchData: MatchWithoutGoals[]
}

const UpcomingMatches: React.FC<UpcomingMatchesProps> = ({matchData}) => {
  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>Upcoming Fixtures</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.tableElement}>Home</th>
            <th className={styles.tableElement}>Match Date</th>
            <th className={styles.tableElement}>Away</th>
          </tr>
        </thead>
        <tbody>
          {matchData.map((match, index) => (
            <tr key={index}>
              <td className={styles.tableElement}>{match.home}</td>
              <td className={styles.tableElement}>{match.match_date}</td>
              <td className={styles.tableElement}>{match.away}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default UpcomingMatches;