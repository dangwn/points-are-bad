import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/router';
import { useMutation } from 'react-query';

import { updateMatch } from '../../lib/adminRequests';
import styles from '../../styles/admin/MatchAdmin.module.css';
import type { MatchWithId } from '../../types/match';

interface MatchAdminProps {
  inputMatchData: MatchWithId[]
}

const baseMatch: MatchWithId = {
  match_id: -1,
  home: '',
  away: '',
  match_date: '',
  home_goals: null,
  away_goals: null
}

const MatchAdmin: React.FC<MatchAdminProps> = ({ inputMatchData }) => {
  const [matchData, setMatchData] = useState<MatchWithId[]>(inputMatchData);
  const [isEditingMatch, setIsEditingMatch] = useState<MatchWithId>(baseMatch)
  const [isEditingIndex, setIsEditingIndex] = useState<number|null>(null);

  const [updateMatchesError, setUpdateMatchesError] = useState<string>('');
  const router = useRouter();

  useEffect(() => {
    if (isEditingIndex !== null && isEditingIndex < matchData.length) {
      setIsEditingMatch(matchData[isEditingIndex]);
    } else {
      setIsEditingMatch(baseMatch)
    }
  }, [isEditingIndex, matchData])

  const submitEditedMatch = useMutation(async () => {
    try {
      const newMatch: MatchWithId = await updateMatch(isEditingMatch);
      setMatchData(
        matchData.map((value) => {
          if (value.match_id === newMatch.match_id) {
            return newMatch;
          } else {
            return value;
          }
        })
      );
      setIsEditingIndex(null);
    } catch {
      setUpdateMatchesError('Could not update match.')
    };
  });

  const preventNegativeInputs = (e: React.KeyboardEvent<HTMLInputElement>): void => {
    if (e.key === '-' || e.key === '_') {
      e.preventDefault();
    };
  };

  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>All Matches</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.tableElement}>Home</th>
            <th className={styles.tableElement}>Home Goals</th>
            <th className={styles.tableElement}>Match Date</th>
            <th className={styles.tableElement}>Away Goals</th>
            <th className={styles.tableElement}>Away</th>
            <th className={styles.tableElement}></th>
          </tr>
        </thead>
        <tbody>
          {matchData.map((match, index) => (
            <React.Fragment key={index}>
            {
              (isEditingIndex===index) ?
              <React.Fragment>
              <tr key={index}>
                <td className={styles.tableElement}>
                  <input 
                    className={styles.tableInput}
                    type='text'
                    placeholder={match.home}
                    onChange={(e) => {setIsEditingMatch(prevState => ({...prevState, home: e.target.value}))}}
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
                    onChange={(e) => {setIsEditingMatch(prevState => ({...prevState, home_goals: parseInt(e.target.value)}))}}
                  />  
                </td>
                <td className={styles.tableElement}>
                  <input 
                    className={styles.tableInput}
                    type='date'
                    placeholder={match.match_date}
                    onChange={(e) => {setIsEditingMatch(prevState => ({...prevState, date: e.target.value}))}}
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
                    onChange={(e) => {setIsEditingMatch(prevState => ({...prevState, away_goals: parseInt(e.target.value)}))}}
                  />  
                </td>
                <td className={styles.tableElement}>
                  <input 
                    className={styles.tableInput}
                    type='text'
                    placeholder={match.away}
                    onChange={(e) => {setIsEditingMatch(prevState => ({...prevState, away: e.target.value}))}}
                  />
                </td>
                <td className={styles.tableElement}>
                  <button className={styles.button} onClick={() => setIsEditingIndex(null)}>Cancel</button>
                </td>
              </tr>
              <tr>
                <td colSpan={6}>
                  <button className={styles.button} onClick={() => submitEditedMatch.mutate()}>Submit Change</button>
                  <button className={styles.button}>Delete Match</button>
                </td>
              </tr>
              </React.Fragment>:
              <tr key={index}>
                <td className={styles.tableElement}>{match.home}</td>
                <td className={styles.tableElement}>{match.home_goals}</td>
                <td className={styles.tableElement}>{match.match_date}</td>
                <td className={styles.tableElement}>{match.away_goals}</td>
                <td className={styles.tableElement}>{match.away}</td>
                <td className={styles.tableElement}>
                  <button className={styles.button} onClick={() => setIsEditingIndex(index)}>Edit</button>
                </td>
              </tr>
            }
            </React.Fragment>
          ))}
        </tbody>
      </table>
      <button className={styles.button}>Create</button>
      <div>{JSON.stringify(isEditingMatch)}</div>
    </div>
  );
};

export default MatchAdmin;