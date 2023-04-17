const accessTokenKey = 'access_token'

export const getAccessToken = (): string|null => {
  return localStorage.getItem(accessTokenKey);
};

export const setAccessToken = (accessToken: string): void => {
  localStorage.setItem(accessTokenKey, accessToken);
}

export const deleteAccessToken = (): void => {
  localStorage.removeItem(accessTokenKey);
}