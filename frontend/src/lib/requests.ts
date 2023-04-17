import { getAccessToken, setAccessToken } from './accessToken';
import { API_HOST } from './constants';

import type { SessionUser } from '../types/user';
import type { SessionUserPoints } from '../types/points';
import type { MatchWithoutGoals } from '../types/match';
import type { LeaderboardApiResponse, LeaderboardUser } from '../types/leaderboard';
import type { Token } from '../types/token';
import type { UserPrediction } from '../types/predictions';

export const createUser = async (token: string, username: string, password: string): Promise<Token> => {
  const response = await fetch(
    `${API_HOST}/user/`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'accept': 'application/json'
      },
      body: JSON.stringify({
        token: token,
        username: username,
        password: password
      }),
      credentials: 'include'
    }
  );
  if (!response.ok) {
    throw new Error('Could not create new user.')
  };
  return response.json();
}

export const getVerificationToken = async (email: string): Promise<Token> => {
  const response = await fetch(
    `${API_HOST}/auth/verify/?email=${email}`,
    {
      method: 'POST',
      headers: {
        'accept': 'application/json'
      }
    }
  );
  if (!response.ok){
    throw new Error('Error creating verification token.');
  };

  return response.json();
}

export const getLeaderboard = async (pageIndex: number, pageSize: number): Promise<LeaderboardApiResponse> => {
  const offset: number = pageIndex * pageSize;
  const response: Response = await fetch(
    `${API_HOST}/points/leaderboard?limit=${pageSize}&offset=${offset}`
  )
  
  if (!response.ok) {
    return {data: [], hasMore: false, pageIndex: pageIndex}
  };
  const data: LeaderboardUser[] = await response.json();
  if (data.length < pageSize) {
    return {data: data, hasMore: false, pageIndex: pageIndex};
  }
  return {data: data, hasMore: true, pageIndex: pageIndex};
};

export const getSessionUser = async (): Promise<SessionUser> => {
  const accessToken: string|null = getAccessToken();
  if (accessToken === null || accessToken === 'undefined') {
    throw new Error('Auth token could not be found');
  };

  const response: Response = await fetch(`${API_HOST}/user/`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });

  if (!response.ok) {
    throw new Error('Error fetching current user');
  };
  return response.json()

}

export const getSessionUserPoints = async (): Promise<SessionUserPoints> => {
  const accessToken: string|null = getAccessToken();
  if (accessToken === null || accessToken === 'undefined') {
    throw new Error('Auth token could not be found');
  };

  const response: Response = await fetch(`${API_HOST}/points/`, {
    headers: {
      Authorization: `Bearer ${accessToken}`,
    },
  });

  if (!response.ok) {
    throw new Error('Error fetching user points');
  };
  return response.json()
};

export const getUpcomingMatches = async (): Promise<MatchWithoutGoals[]> => {
  const today: Date = new Date();
  const todayDateString: string = today.toISOString().slice(0, 10);

  const response: Response = await fetch(`${API_HOST}/match/?start_date=${todayDateString}`);
  if (!response.ok){
    throw new Error('Error fetching upcoming matches');
  };
  return response.json();
};

export const getUpcomingUserPredictions = async(): Promise<UserPrediction[]> => {
  const accessToken: string|null = getAccessToken();
  if (accessToken === null || accessToken === 'undefined') {
    throw new Error('Auth token could not be found');
  };
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
    throw new Error('Error fetching predictions data');
  };
  return response.json();
}

export const logUserIn = async (email: string, password: string): Promise<Token> => {
  const response = await fetch(
    `${API_HOST}/auth/login/`,
    {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'accept': 'application/json'
      },
      body: JSON.stringify({
        email: email,
        password: password
      }),
      credentials: 'include'
    }
  );

  if (!response.ok) {
    throw new Error('Could not log in');
  };

  return response.json();
}

export const refreshAccessToken = async (): Promise<Token> => {
  const response: Response = await fetch(
    `${API_HOST}/auth/refresh/`,
    {
      method: 'POST',
      headers: {
        'accept': 'application/json'
      },
      credentials: 'include'
    }
  );
  if (!response.ok) {
    throw new Error('Error fetching refreshed access token');
  };
  return response.json();
};