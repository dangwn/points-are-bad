import React from 'react';

import { NULL_PREDICTION_PENALTY } from '@/lib/constants';
import styles from '@/styles/PredictionsPage.module.css'
import type { UserPrediction } from '@/types/predictions';

interface ResultsTableResults {
  predictions: UserPrediction[]
}

const calculateMatchPoints = (pred: UserPrediction): number => {
  if (pred.match.home_goals === null || pred.match.away_goals === null){
    return 0;
  }
  else if (pred.home_goals === null || pred.away_goals === null) {
    return NULL_PREDICTION_PENALTY;
  }
  return Math.abs((pred.home_goals - pred.match.home_goals) + (pred.away_goals - pred.match.away_goals));
}

const displayScore = (goals: number|null): string|number => {
  if (goals === null) {
    return '-';
  }
  return goals;
}

const ResultsTable: React.FC<ResultsTableResults> = ({ predictions }) => {
  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>Results</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}></th>
            <th className={styles.th}>Predicted</th>
            <th className={styles.th}>Score</th>
            <th className={styles.th}>Points</th>
          </tr>
        </thead>
        <tbody>
          {
            predictions.map((prediction, index) => (
              <>
                <tr>
                  <td colSpan={4}>
                    <hr className={styles.hr}/>
                  </td>
                </tr>
                <tr key={index * 2}>
                  <td className={styles.tableElement}>{prediction.match.home}</td>
                  <td className={styles.tableElement}>{displayScore(prediction.home_goals)}</td>
                  <td className={styles.tableElement}>{displayScore(prediction.match.home_goals)}</td>
                  <td className={styles.tableElement} rowSpan={2}>{calculateMatchPoints(prediction)}</td>
                </tr>
                <tr key={index * 2 + 1}>
                  <td className={styles.tableElement}>{prediction.match.away}</td>
                  <td className={styles.tableElement}>{displayScore(prediction.away_goals)}</td>
                  <td className={styles.tableElement}>{displayScore(prediction.match.away_goals)}</td>
                </tr>
              </>
            ))
          }
        </tbody>
      </table>
    </div>
  )
}

export default ResultsTable;