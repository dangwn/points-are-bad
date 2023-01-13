import axios from "axios";
import { APIHOST } from "../constants";

export const getCurrentUserInfo = async (accessToken) => {
  return await axios.get(
    APIHOST+'user/',
    {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${accessToken}`
      }
    }
  )
    .then(({ data }) => {
      return {'username': data['username'], is_admin: data['is_admin']};
    })
    .catch(err => {
      console.log('Error trying to retrieve user information', err);
    });
};

export const getCurrentUserScore = async (accessToken) => {
  return await axios.get(
    APIHOST+'points/',
    {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${accessToken}`
      }
    }
  )
    .then(({ data }) => {
      return data;
    })
    .catch(err => {
      console.log('Error trying to retrieve user points', err);
    });
};