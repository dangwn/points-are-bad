import React from "react";
import { useQuery } from "react-query";

import styles from "../../styles/index/IndexTables.module.css";
import { API_HOST, QUERY_OPTIONS } from "../../lib/constants";

type Match = {
  match_date: string,
  home: string, 
  away: string
}

const UpcomingMatches: React.FC = () => {
  const { data, isLoading, isError } = useQuery<Match[]>(
    'upcomingMatches',
    async () => {
      const today = new Date();
      const todayDateString = today.toISOString().slice(0, 10);

      const response = await fetch(`${API_HOST}/match/?start_date=${todayDateString}`);
    
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

  const matchData: Match[] = data || []
  
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