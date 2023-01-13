import axios from "axios";
import Cookies from 'js-cookie';

import { APIHOST } from "../constants";

export const loginRequest = async (username, password) => {
  return await axios.post(
    APIHOST+'login/',
    new URLSearchParams({
      'grant_type': '',
      'username': username,
      'password': password,
      'scope': '',
      'client_id': '',
      'client_secret': ''
    }),
    {
      headers: {
        'accept': 'application/json'
      },
      withCredentials: true
    }
  )
    .then(({ data }) => {
      return data;
    })
    .catch(err => {
      console.log('Error trying to log in:', err);
    });
};

export const createUserRequest = async (new_username, new_email, new_password) => {   
  return await axios.post(
    APIHOST+'user/',
    {
      username: new_username,
      email: new_email,
      password: new_password
    },
    {
      headers: {
        'accept': 'application/json',
        'Content-Type': 'application/json'
      },
      withCredentials: true
    }
  )
    .then(({ data }) => {
      return data;
    })
    .catch(err => {
      console.log('Error trying to create user:', err);
    });
};

export const getRefreshedAccessToken = async () => {
  if (Cookies.get('csrf_refresh_token') == undefined) {
    return 
  }

  return await axios.post(
    APIHOST+'login/refresh/',
    {},
    {
      headers: {
        'accept': 'application/json',
        'X-CSRF-TOKEN': Cookies.get('csrf_refresh_token')
      },
      withCredentials: true
    }
  )
    .then(({ data }) => {
      return data
    })
    .catch(err => {
      console.log('Error trying to get refreshed token:', err);
    });
}

export const logoutRequest = async (accessToken) => {
  return await axios.delete(
    APIHOST+'login/',
    {
      headers: {
        'accept': 'application/json',
        'Authorization': `Bearer ${accessToken}`
      },
      withCredentials: true
    }
  ).catch(err => {
    console.log('Could not logout due to:', err);
  }) 
}