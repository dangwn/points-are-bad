const accessTokenKey = 'access_token'

export const getAccessToken = (): string => {
  const accessToken: string|null = localStorage.getItem(accessTokenKey);
  if (accessToken === null || accessToken === 'undefined') {
    throw new Error('Auth token could not be found');
  };
  return accessToken;
};

export const setAccessToken = (accessToken: string): void => {
  localStorage.setItem(accessTokenKey, accessToken);
}

export const deleteAccessToken = (): void => {
  localStorage.removeItem(accessTokenKey);
}