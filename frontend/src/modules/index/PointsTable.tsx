import React from 'react';

import Cell from './Cell';
import { createPositionString } from '@/lib/change';

import styles from '@/styles/HomePage.module.css';

interface PointsTableProps {
  username: string,
  points: number,
  correctScores: number,
  largestError: number,
  position: number|null
}

const PointsTable: React.FC<PointsTableProps> = ({username, points, correctScores, largestError, position}) => {
   return (
    <div className={styles.pointsTableContainer}>
    <div className={styles.pointsTable}>
      <Cell 
        className={styles.topCell}
        topText={username} 
        bottomText={createPositionString(position)} 
        topSize='3em' 
        bottomSize='1.5em'
        />
      <div className={styles.pointsContainer}>
        <Cell 
          className={styles.cell}
          topText={correctScores} 
          bottomText='Correct Scores' 
          topSize='2em' 
          bottomSize='1em'
        />
        <Cell 
          className={styles.cell}
          topText={points} 
          bottomText='Points' 
          topSize='2.6em' 
          bottomSize='1.3em'
        />
        <Cell 
          className={styles.cell}
          topText={largestError} 
          bottomText='Largest Error' 
          topSize='2em' 
          bottomSize='1em'
        />
      </div>
    </div></div>
  );
};

export default PointsTable;
