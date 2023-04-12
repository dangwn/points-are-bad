import React, { useState } from "react";
import { useQuery } from "react-query";

import styles from "../../styles/index/IndexTables.module.css";
import { API_HOST, QUERY_OPTIONS } from "../../lib/constants";

type userPrediction =   {
  prediction_id: number,
  home_goals: number|undefined,
  away_goals: number|undefined,
  match: {
    match_date: string,
    home: string,
    away: string,
    home_goals: number|undefined,
    away_goals: number|undefined
  },
  user: {
    username: string,
    is_admin: boolean
  }
}

const PredictionsTable: React.FC = () => {
  const [editedPredictions, setEditedPredictions] = useState<userPrediction[]>([])

  const { data, isLoading, isError } = useQuery<userPrediction[]>(
    'userPredictions',
    async () => {
      const accessToken = localStorage.getItem('access_token');
      const today = new Date();
      const todayDateString = today.toISOString().slice(0, 10);

      const response = await fetch(
        `${API_HOST}/prediction/?start_date=${todayDateString}`,
        {
            headers: {
            'Authorization': `Bearer ${accessToken}`
          }
        }
      );
    
      if (!response.ok){
        throw new Error('Error fetching matches data');
      };
      return response.json();
    },
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <div>Loading...</div>;
  };

  if (isError) {
    return <div>Error fetching upcoming matches</div>
  };

  const predictionsData: userPrediction[] = data || []
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (!/^\d+$/.test(e.key) && e.key !== 'Backspace') {
      e.preventDefault();
    }
  }
  const handleEditedPredictions = (pred: userPrediction) => {
    const newPredictions = editedPredictions;
    for (let newPred of newPredictions) {
      if (newPred.prediction_id === pred.prediction_id) {
        newPred.home_goals = pred.home_goals;
        newPred.away_goals = pred.away_goals;

        setEditedPredictions(newPredictions);
        return;
      }
    }

    newPredictions.push(pred);
    setEditedPredictions(newPredictions);
  }
  
  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>Upcoming Fixtures</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.tableElement}>Home</th>
            <th className={styles.tableElement}>Home Goals</th>
            <th className={styles.tableElement}>Match Date</th>
            <th className={styles.tableElement}>Away Goals</th>
            <th className={styles.tableElement}>Away</th>
          </tr>
        </thead>
        <tbody>
          {predictionsData.map((prediction, index) => (
            <tr key={index}>
              <td className={styles.tableElement}>{prediction.match.home}</td>
              <td 
                className={styles.tableElement} 
                contentEditable
                onKeyDown={handleKeyDown}
                onBlur={
                  (e) => {
                    const newHomeGoals = Number(e.target.textContent);
                    if (!isNaN(newHomeGoals)){
                      predictionsData[index].home_goals = newHomeGoals;
                      const newPred: userPrediction = predictionsData[index];
                      console.log(newPred);
                      handleEditedPredictions(newPred);

                    }
                  }
                }
              >
                {prediction.home_goals ? prediction.home_goals : 'Fill Here'}
              </td>
              <td className={styles.tableElement}>{prediction.match.match_date}</td>
              <td className={styles.tableElement}>
                {prediction.away_goals ? prediction.away_goals : 'Fill Here'}
              </td>
              <td className={styles.tableElement}>{prediction.match.away}</td>
            </tr>
          ))}
          <tr>
            <td colSpan={5}>
              <button onClick={() => {
                console.log(editedPredictions);
              }}>Submit</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}

export default PredictionsTable;