import React from 'react';

import { preventNegativeInputs } from '@/lib/change';
import styles from '@/styles/admin/MatchAdmin.module.css'
import type { MatchWithId } from '@/types/match';

interface MatchEditingRowProps {
  match: MatchWithId;
  index: number;
  homeOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  homeGoalsOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  matchDateOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  awayGoalsOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  awayOnChange: (e: React.ChangeEvent<HTMLInputElement>) => void;
  cancelOnClick: () => void;
  submitOnClick: () => void;
  deleteOnClick: () => void
}

const MatchEditingRow: React.FC<MatchEditingRowProps> = ({
  match,
  index,
  homeOnChange,
  homeGoalsOnChange,
  matchDateOnChange,
  awayGoalsOnChange,
  awayOnChange,
  cancelOnClick,
  submitOnClick,
  deleteOnClick
}) => {
  return (
    <>
      <tr key={index}>
        <td className={styles.tableElement}>
          <input 
            className={styles.tableInput}
            type='text'
            placeholder={match.home}
            onChange={homeOnChange}
          />
        </td>
        <td className={styles.tableElement}>
          <input 
            className={styles.tableInput}
            type='number'
            onKeyDown={preventNegativeInputs}
            min={0}
            placeholder={
              (match.home_goals !== null && !isNaN(match.home_goals)) ?
              match.home_goals.toString() :
              'Home goals'
            }
            onChange={homeGoalsOnChange}
          />  
        </td>
        <td className={styles.tableElement}>
          <input 
            className={styles.tableInput}
            type='date'
            placeholder={match.match_date}
            onChange={matchDateOnChange}
          />
        </td>
        <td className={styles.tableElement}>
          <input 
            className={styles.tableInput}
            type='number'
            onKeyDown={preventNegativeInputs}
            min={0}
            placeholder={
              (match.away_goals !== null && !isNaN(match.away_goals)) ?
              match.away_goals.toString() :
              'Away goals'
            }
            onChange={awayGoalsOnChange}
          />  
        </td>
        <td className={styles.tableElement}>
          <input 
            className={styles.tableInput}
            type='text'
            placeholder={match.away}
            onChange={awayOnChange}
          />
        </td>
        <td className={styles.tableElement}>
          <button className={styles.button} onClick={cancelOnClick}>Cancel</button>
        </td>
      </tr>
      <tr>
        <td colSpan={6}>
          <button className={styles.button} onClick={submitOnClick}>Submit Change</button>
          <button className={styles.button} onClick={deleteOnClick}>Delete Match</button>
        </td>
      </tr>
    </>
  )
}

export default MatchEditingRow;