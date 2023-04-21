import { getAccessToken, setAccessToken } from './accessToken';
import { API_HOST } from './constants';

import type { SessionUser } from '../types/user';
import type { LeaderboardPoints } from '../types/points';
import type { MatchWithoutGoals, Match } from '../types/match';
import type { LeaderboardApiResponse, LeaderboardUser } from '../types/leaderboard';
import type { Token } from '../types/token';
import type { UserPrediction, NewPrediction } from '../types/predictions';

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
  const accessToken: string = getAccessToken();

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

export const getSessionUserPoints = async (): Promise<LeaderboardPoints> => {
  const accessToken: string = getAccessToken();

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
  const accessToken: string = getAccessToken();
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
};

export const getMatchesWithGoals = async (startDate?: string, endDate?: string): Promise<Match> => {
  const queryParams: any = {}
  if (startDate) {
    queryParams['start_date'] = startDate;
  };
  if (endDate) {
    queryParams['end_date'] = endDate;
  };
  const queryParamsString: string = new URLSearchParams(queryParams).toString();
  
  const response: Response = await fetch(`${API_HOST}/match/with-goals/?${queryParamsString}`);
  if (!response.ok){
    throw new Error('Error fetching upcoming matches');
  };
  return response.json();
};

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

export const updateUsername = async (newUsername: string) => {
  const accessToken: string = getAccessToken();
  const response = await fetch(`${API_HOST}/user/username`,
   {
      method: 'PUT',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'accept': 'application/json'
      },
      body: newUsername,
      credentials: 'include'
   }
  );
  if (!response.ok){
    throw new Error('Could not update username');
  };
  return response.json()
}

export const updateUserPredictions = async (newUserPredictions: NewPrediction[]): Promise<void> => {
  const accessToken: string = getAccessToken();
  const requestBody: string = JSON.stringify(newUserPredictions);

  const response = await fetch(`${API_HOST}/prediction/`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json'
    },
    body: requestBody,
    credentials: 'include',
  });

  if (!response.ok) {
    throw new Error('Could not update user predictions')
  };
}