import axios from "axios";

import { API_HOST } from "../constants";

import type UserPointsData from "@/types/UserPointsData";

export const getGlobalLeaderBoard = async (
  tableStart: number,
  tableEnd: number
): Promise<UserPointsData[]> => {
  return await axios.get(
    API_HOST+'points/leaderboard',
    {
      params: {
        'table_start': tableStart,
        'table_end': tableEnd
      },
      headers: {
        'accept': 'application/json'
      }
    }
  ).then(({ data }) => {
    return data;
  }).catch((err: Error) => {
    console.log(err)
  });
} 