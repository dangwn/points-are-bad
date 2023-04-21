import React, { useState } from 'react';
import { useMutation } from 'react-query';

import { setAccessToken } from '../../lib/accessToken';
import { logUserIn } from '../../lib/requests';
import styles from '@/styles/LoginPage.module.css'
import type { Token } from '../../types/token';
import type { LoginData } from '../../types/auth';

type LoginFormProps = {
  onSuccess: () => void;
};

const LoginForm: React.FC<LoginFormProps> = ({ onSuccess }) => {
  const [formData, setFormData] = useState<LoginData>({
    email: '',
    password: '',
  });
  const [logInError, setLogInError] = useState<string>('');

  const mutation = useMutation(
    async () => {
      const token: Token = await logUserIn(
        formData.email,
        formData.password
      );

      setAccessToken(token.access_token);
    },
    {
      onSuccess,
    }
  );

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setLogInError('');
    mutation.mutate();
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prevData) => ({
      ...prevData,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div>
      <form className={styles.loginForm} onSubmit={handleSubmit}>
        <input
          className={styles.input}
          type='email'
          name='email'
          value={formData.email}
          onChange={handleChange}
          placeholder='Email'
        />
        <input
          className={styles.input}
          type='password'
          name='password'
          value={formData.password}
          onChange={handleChange}
          placeholder='Password'
        />
        <button className={styles.button} type='submit' disabled={mutation.isLoading}>
          {mutation.isLoading ? 'Loading...' : 'Sign In'}
        </button>
      </form>
      {(logInError === '') ? null : (<div className={styles.error}>{logInError}</div>)}
    </div>
  );
};

export default LoginForm;
