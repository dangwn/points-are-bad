import { useState } from 'react';
import { useAuth } from '../hooks';
import { useNavigate } from 'react-router-dom';

import { loginRequest } from '../requests/authRequests';

export const LogInForm = () => {
  const [loginStatus, setLoginStatus] = useState(null);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');

  const {setAccessToken} = useAuth('')

  const navigate = useNavigate();
  const routeChange = () => {
    let homePath = '/';
    navigate(homePath);
  };

  const handleLogin = async (event) => {
    event.preventDefault();
    
    let newAccessToken = null;
    try {
      const accessTokenResponse = await loginRequest(username, password);

      if (accessTokenResponse == undefined) {
        throw new Error('Unable to login with given credentials');
      };
      newAccessToken = accessTokenResponse['access_token'];
    } catch (err) {
      setLoginStatus('Incorrect username or password');
      console.log(err);
    };

    if (newAccessToken != null) {
      console.log('Setting access token and rerouting to home page');
      setAccessToken(newAccessToken);
      routeChange();
    };
  };

  return (
    <>
      <form onSubmit={handleLogin}>
        <h3>Points Are Bad</h3>
        <label>
          <input
            type='text'
            placeholder='Username'
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
        <div></div>
        <label>
          <input
            type='password'
            placeholder='Password'
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        <div></div>
        <button type='submit'>Login</button>
      </form>
      <div>{loginStatus}</div>
    </>
  );
};
