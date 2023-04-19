import { getAccessToken } from './accessToken';
import { API_HOST } from './constants';
import type { MatchWithId, MatchWithoutGoals } from '../types/match';

export const getFullMatches = async (startDate?: string, endDate?: string): Promise<MatchWithId[]> => {
  const accessToken: string = getAccessToken();
  const queryParams: any = {}
  if (startDate) {
    queryParams['start_date'] = startDate;
  };
  if (endDate) {
    queryParams['end_date'] = endDate;
  };
  const queryParamsString: string = new URLSearchParams(queryParams).toString();

  const response: Response = await fetch(
    `${API_HOST}/match/full/?${queryParamsString}`, 
    {
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'Content-Type': 'application/json'
      },
      credentials: 'include'
    }
  );

  if (!response.ok){
    throw new Error('Error fetching upcoming matches');
  };

  return response.json();
}

export const createMatch = async (match: MatchWithoutGoals): Promise<MatchWithId> => {
  const accessToken: string = getAccessToken();
  const requestBody: string = JSON.stringify(match);

  const response: Response = await fetch(`${API_HOST}/match/`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
      'accept': 'application/json'
    },
    body: requestBody,
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Could not create match');
  };

  return response.json();
}

export const updateMatch = async (match: MatchWithId): Promise<MatchWithId> => {
  const accessToken: string = getAccessToken();
  const requestBody: string = JSON.stringify(match);

  const response: Response = await fetch(`${API_HOST}/match/`, {
    method: 'PUT',
    headers: {
      'Authorization': `Bearer ${accessToken}`,
      'Content-Type': 'application/json',
    },
    body: requestBody,
    credentials: 'include',
  });
  if (!response.ok) {
    throw new Error('Could not update matches');
  };

  return response.json();
}

export const deleteMatch = async (matchId: number): Promise<void> => {
  const accessToken: string = getAccessToken();

  const response: Response = await fetch(
    `${API_HOST}/match/${matchId}`,
    {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${accessToken}`,
        'accept': '*/*'
      },
      credentials: 'include'
    }
  );

  if (!response.ok) {
    throw new Error('Could not delete match');
  };
}

