import { createContext, useContext, useState } from "react";
import Cookies from "js-cookie";

import { getRefreshedAccessToken } from './requests/authRequests';

export const AuthContext = createContext(null);
export const useAuth = () => useContext(AuthContext);

export const AuthProvider = ({children}) => {
  const [accessToken, setAccessToken] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  return (
    <AuthContext.Provider value={{
      accessToken, setAccessToken, 
      isLoading, setIsLoading
    }}>
      {children}
    </AuthContext.Provider>
  )
};

export const refreshAccessToken = () => {
  const { setAccessToken } = useAuth();

  if (Cookies.get('csrf_refresh_token') == null) {
    setAccessToken(null);
  } else {
    getRefreshedAccessToken().then(({ access_token }) => {
      setAccessToken(access_token || null);
    }).catch(err => console.log(err))
  };
}
