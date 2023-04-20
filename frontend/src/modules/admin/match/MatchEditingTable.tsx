import React, { useState, useEffect } from 'react';
import { useMutation } from 'react-query';

import MatchEditingRow from './MatchEditingRow';
import { updateMatch, deleteMatch } from '@/lib/adminRequests';
import styles from '@/styles/admin/MatchAdmin.module.css';
import type { MatchWithId } from '@/types/match';

interface MatchEditingTableProps {
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

const MatchEditingTable: React.FC<MatchEditingTableProps> = ({ inputMatchData }) => {
  const [matchData, setMatchData] = useState<MatchWithId[]>(inputMatchData);
  const [isEditingMatch, setIsEditingMatch] = useState<MatchWithId>(baseMatch)
  const [isEditingIndex, setIsEditingIndex] = useState<number|null>(null);
  const [updateMatchesError, setUpdateMatchesError] = useState<string>('');

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
      setUpdateMatchesError('');
    } catch {
      setUpdateMatchesError('Could not update match.');
    };
  });

  const submitDeleteMatch = useMutation(async () => {
    try {
      await deleteMatch(isEditingMatch.match_id);
      setMatchData(matchData.filter((match) => (match.match_id !== isEditingMatch.match_id)));
      setIsEditingIndex(null);
      setUpdateMatchesError('');
    } catch {
      setUpdateMatchesError('Could not delete match.');
    }
  });

  return (
    <div>
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
            <>
            {
              (isEditingIndex===index) ?
              <MatchEditingRow
                match={match}
                index={index}
                homeOnChange={(e) => {setIsEditingMatch(prevState => ({...prevState, home: e.target.value}))}}
                homeGoalsOnChange={(e) => {setIsEditingMatch(prevState => ({...prevState, home_goals: parseInt(e.target.value, 10)}))}}
                matchDateOnChange={(e) => {setIsEditingMatch(prevState => ({...prevState, date: e.target.value}))}}
                awayGoalsOnChange={(e) => {setIsEditingMatch(prevState => ({...prevState, away_goals: parseInt(e.target.value, 10)}))}}
                awayOnChange={(e) => {setIsEditingMatch(prevState => ({...prevState, away: e.target.value}))}}
                cancelOnClick={() => setIsEditingIndex(null)}
                submitOnClick={() => submitEditedMatch.mutate()}
                deleteOnClick={() => submitDeleteMatch.mutate()}
              />
              :
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
            </>
          ))}
        </tbody>
      </table>
      <div className={styles.Error}>{updateMatchesError}</div>
    </div>
  );
};

export default MatchEditingTable;