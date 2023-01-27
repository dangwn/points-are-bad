import axios from "axios";

import { API_HOST } from "../constants";

interface userPoints {
  points: number,
  correct_scores: number,
  largest_error: number,
  position: number
}

export const getUserPoints = async (accessToken: string, provider: string): Promise<userPoints> => {
  return await axios.get(
    API_HOST+'points/position/',
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'provider': provider,
        'accept': 'application/json'
        }
    }
  ).then(({ data }) =>{
    return data;
  }).catch(() => {
    return {
      points: 0,
      correct_scores: 0,
      largest_error: 0,
      position: 0
    }
  });
}