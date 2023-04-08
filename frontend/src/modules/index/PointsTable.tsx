import React from "react";

import styles from "../../styles/index/IndexTables.module.css";

interface PointsTableProps {
  username: string,
  points: number,
  correctScores: number,
  largestError: number
}

const PointsTable: React.FC<PointsTableProps> = ({username, points, correctScores, largestError}) => {
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
