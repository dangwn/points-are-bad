import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import { createUserRequest } from '../requests/authRequests';
import { useAuth } from '../hooks';

export const SignUpForm = () => {
  const [userCreationStatus, setUserCreationStatus] = useState(null);
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [email, setEmail] = useState('');
  const { setAccessToken } = useAuth();

  const navigate = useNavigate();
  const routeChange = () => {
    let homePath = '/';
    navigate(homePath);
  };

  const handleSumbitNewUser = async (event) => {
    event.preventDefault();

    let newAccessToken = null;
    try {
      const accessTokenResponse = await createUserRequest(username, email, password);

      if (accessTokenResponse == undefined) {
        throw new Error('Unable to sign up with given information');
      };
      newAccessToken = accessTokenResponse['access_token'];
    } catch (err) {
      setUserCreationStatus('Username or email not available');
      console.log(err);
    };

    if (newAccessToken != null) {
      console.log('Setting access token and rerouting to home page');
      setAccessToken(newAccessToken);
      routeChange();
    };
  }

  return (
    <>
      <form onSubmit={handleSumbitNewUser}>
        <h3>Sign Up</h3>
        <div>
          <label>  
            <input
              type='text'
              placeholder='Username'
              onChange={(e) => setUsername(e.target.value)}
            />
          </label>
        </div>
        <div>
          <label>
            <input
              type='text'
              placeholder='Email Address'
              onChange={(e) => setEmail(e.target.value)}
            />
          </label>
        </div>
        <div>
          <label>
            <input
              type='password'
              placeholder='Password'
              onChange={(e) => setPassword(e.target.value)}
            />
          </label>
        </div>
        <div>
          <button type='submit'>Submit</button>
        </div>
      </form>
      <p>{userCreationStatus}</p>
    </>
  )
};
