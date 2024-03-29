import React, { useState } from 'react';
import { useMutation } from 'react-query';
import { useRouter } from 'next/router';

import { updateUserPredictions } from '@/lib/requests';
import { preventNegativeInputs } from '@/lib/change';
import styles from '@/styles/PredictionsPage.module.css';
import type { UserPrediction, NewPrediction } from '@/types/predictions';

interface PredictionTableProps {
  predictions: UserPrediction[]
}

const PredictionsTable: React.FC<PredictionTableProps> = ({predictions}) => {
  const [predictionsData, setPredictionsData] = useState<UserPrediction[]>(predictions);
  const [isEditing, setIsEditing] = useState<boolean>(false);
  const [predictionError, setPredictionError] = useState<string>('');
  const router = useRouter();

  const submitEditedPredictions = useMutation(async () => {
    // Reduce data to new prediction type, and remove nulls and NaNs
    const newUserPredictions: NewPrediction[] = predictionsData.map((pred) => {
        return {
          prediction_id: pred.prediction_id,
          home_goals: pred.home_goals,
          away_goals: pred.away_goals
        }
      }).filter(
        (pred) => (
              pred.home_goals !== null && 
              !isNaN(pred.home_goals) && 
              pred.away_goals !== null && 
              !isNaN(pred.away_goals)
            )
      );
      // No need to send the request if there are no predictions to send
      if (newUserPredictions.length === 0) {
        router.reload();
      };

    try {
      await updateUserPredictions(newUserPredictions);
      router.reload();
    } catch {
      setPredictionError("Something went wrong when updating your predictions. Please try again.");
    }
  });

  const handlePredictionsChange = (
    predictionId: number,
    goals: number|null,
    isHome: boolean
  ) => {
    const newPredictions: UserPrediction[] = [...predictionsData]
    const predictionIndex: number = newPredictions.findIndex(
      pred => (pred.prediction_id === predictionId)
    );
    if (isHome) {
      newPredictions[predictionIndex].home_goals = goals
    } else {
      newPredictions[predictionIndex].away_goals = goals
    };
    setPredictionsData(newPredictions);
  };

  return (
    <div className={styles.container}>
      <div className={styles.tableTitle}>Predictions</div>
      <table className={styles.table}>
        <thead>
          <tr>
            <th className={styles.th}>Home</th>
            <th className={styles.th}>Home Goals</th>
            <th className={styles.th}>Match Date</th>
            <th className={styles.th}>Away Goals</th>
            <th className={styles.th}>Away</th>
          </tr>
        </thead>
        <tbody>
          {predictionsData.map((prediction, index) => (
            <tr key={index}>
              <td className={styles.tableElement}>{prediction.match.home}</td>
              { 
                isEditing ? 
                  <input 
                    className={styles.tableInput}
                    type='number'
                    onKeyDown={preventNegativeInputs}
                    min={0}
                    placeholder={
                      (prediction.home_goals !== null && !isNaN(prediction.home_goals)) ?
                      prediction.home_goals.toString() :
                      'Home goals'
                    }
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      handlePredictionsChange(
                        prediction.prediction_id,
                        parseInt(e.target.value, 10),
                        true
                      );
                    }}
                  /> :
                  <td className={styles.tableElement}>
                    {(prediction.home_goals !== null) ? prediction.home_goals : '-'}
                  </td>
                }
              <td className={styles.tableElement}>{prediction.match.match_date}</td>
              { 
                isEditing ? 
                  <input 
                    className={styles.tableInput}
                    type='number'
                    onKeyDown={preventNegativeInputs}
                    min={0}
                    placeholder={
                      (prediction.away_goals !== null && !isNaN(prediction.away_goals)) ? 
                      prediction.away_goals.toString() : 
                      'Away goals'
                    }
                    onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                      handlePredictionsChange(
                        prediction.prediction_id,
                        parseInt(e.target.value, 10),
                        false
                      );
                    }}
                  /> :
                  <td 
                    className={styles.tableElement}
                  >
                    {(prediction.away_goals !== null) ? prediction.away_goals : '-'}
                  </td>
                }
              <td className={styles.tableElement}>{prediction.match.away}</td>
            </tr>
          ))}
          {
            isEditing ? 
            <tr>
              <td colSpan={5}>
                <button className={styles.button} onClick={() => {
                  submitEditedPredictions.mutate();
                }}>Save</button>
              </td>
            </tr> :
            null 
          }
          <tr>
            <td colSpan={5}>
              <button className={styles.button} onClick={() => {
                setIsEditing(!isEditing);
                setPredictionError('');
              }}>
                {isEditing? 'Cancel' : 'Edit'} 
              </button>
            </td>
          </tr>
        </tbody>
      </table>
      {
        (predictionError === '') ? 
        null :
        <div className={styles.predictionError}>{predictionError}</div>
      }
    </div>
  );
}

export default PredictionsTable;