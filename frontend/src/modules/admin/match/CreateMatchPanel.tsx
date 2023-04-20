import React, { useState } from 'react';
import { useMutation } from 'react-query';

import { createMatch } from '@/lib/adminRequests';
import styles from '@/styles/admin/MatchAdmin.module.css';
import type { MatchWithoutGoals, MatchWithId } from '@/types/match';

interface CreateMatchPanelProps {}

const CreateMatchPanel: React.FC<CreateMatchPanelProps> = ({}) => {
  const [newMatch, setNewMatch] = useState<MatchWithoutGoals>({
    match_date:'',
    home: '',
    away: ''
  });
  const [isCreatingNewMatch, setIsCreatingNewMatch] = useState<boolean>(false)
  const [newMatchMessage, setNewMatchMessage] = useState<string>('');
  
  const handleCreateMatch = useMutation(async () => {
    try {
      const match: MatchWithId = await createMatch(newMatch);
      if (!match) {
        throw new Error;
      }
      setNewMatchMessage('Match created successfully. Please refresh to see updated table.');
      setIsCreatingNewMatch(false);
    } catch {
      setNewMatchMessage('Could not create new match.');
    }
  });
  
  return (
    <div>
      {
        isCreatingNewMatch ?
        <>
        <div className={styles.newMatchContainer}>
          <input 
            className={styles.newMatchInput}
            type='text' 
            placeholder='Home'
            onChange={(e) => setNewMatch(prevState => ({...prevState, home: e.target.value}))}
          />
          <input 
            className={styles.newMatchInput}
            type='date' 
            placeholder='Match Date'
            onChange={(e) => setNewMatch(prevState => ({...prevState, match_date: e.target.value}))}
          />
          <input 
            className={styles.newMatchInput}
            type='text' 
            placeholder='Away'
            onChange={(e) => setNewMatch(prevState => ({...prevState, away: e.target.value}))}
          />
          </div>
          <div className={styles.newMatchContainer}>
            <button className={styles.newMatchButton} onClick={() => handleCreateMatch.mutate()}>
              Submit
            </button>
            <button className={styles.newMatchButton} onClick={() => setIsCreatingNewMatch(false)}>
              Cancel
            </button>
          </div> 
        </>:
        <button className={styles.button} onClick={() => setIsCreatingNewMatch(true)}>
          Create Match
        </button>
      }
      <h3 className={styles.Error}>{newMatchMessage}</h3>
    </div>
  )
}

export default CreateMatchPanel;