import { useEffect } from "react";
import { Outlet } from "react-router-dom";
import Cookies from "js-cookie";

import { useAuth } from "../hooks";
import { getRefreshedAccessToken } from "../requests/authRequests";

export const LoadingScreen = () => {
  const { setAccessToken, isLoading, setIsLoading } = useAuth();

  useEffect(() => {
    const initToken = async () => {
      let newToken;
      if (Cookies.get('csrf_refresh_token') != null) {
        let resp = await getRefreshedAccessToken();
        newToken = resp['access_token'];
      } 

      setAccessToken(newToken);
      setIsLoading(false);
    };

    initToken();
  }, [isLoading == true]);

  if (isLoading) {
    return <div>Loading</div>
  } else {
    return <Outlet />
  };
};