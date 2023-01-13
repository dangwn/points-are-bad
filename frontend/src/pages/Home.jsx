import { useState, useEffect } from "react";

import { useAuth } from "../hooks";
import { getCurrentUserScore, getCurrentUserInfo } from "../requests/userRequests";
import { logoutRequest } from "../requests/authRequests";

export const Home = () => {
  const [homeData, setHomeData] = useState('');
  const { accessToken, setIsLoading } = useAuth();

  useEffect(() => {
    const getHomePageData = async () => {
      if (accessToken == null) {
        setHomeData('Welcome to points are bad');
      } else {
        const userInfo = await getCurrentUserInfo(accessToken);
        const userScore = await getCurrentUserScore(accessToken);

        setHomeData(
          <>
            <p>{JSON.stringify(userInfo)}</p>
            <p>{JSON.stringify(userScore)}</p>
          </>
        )

      }
    }

    getHomePageData();
  }, []);

  return (
    <>
      <div>
        {homeData}
      </div>
      <button onClick={async () => {
        await logoutRequest(accessToken);
        setIsLoading(true);
      }}>Logout</button>
    </>
  );
};