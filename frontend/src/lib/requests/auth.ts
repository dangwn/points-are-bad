import axios from "axios";

import { API_HOST } from "../constants";

interface userData {
  displayName: string,
  isAdmin: boolean
}

interface newUserRequiredData {
  display_name: string,
  avatar: string
}

export const startUserSession = async (accessToken: string, provider: string): Promise<userData> => {
  return await axios.get(
    API_HOST+'user/',
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'provider': provider,
        'accept': 'application/json'
      }
    }
  )
    .then(({ data }) => {
      return data;
    }).catch((err: Error) => {
      return {
        displayName: '',
        isAdmin: false
      }
    });
};

export const createNewUser = async (
  accessToken: string,
  provider: string,
  userData: newUserRequiredData
): Promise<userData> => {
  return await axios.post(
    API_HOST+'user/',
    userData,
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'provider': provider,
        'accept': 'application/json',
        'Content-Type': 'application/json'
      }
    }
  ).then(({ data }) => {
    return data;
  }).catch((err: Error) => {
    return {
      displayName: '',
      isAdmin: false
    }
  });
}