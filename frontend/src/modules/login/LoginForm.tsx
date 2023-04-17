import React, { useState } from 'react';
import { useMutation } from 'react-query';

import { setAccessToken } from '../../lib/accessToken';
import { logUserIn } from '../../lib/requests';
import styles from '../../styles/login/LogInForm.module.css'
import type { Token } from '../../types/token';
import { LoginData } from '../../types/auth';

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
    <div className={styles.container}>
      <form onSubmit={handleSubmit}>
        <label className={styles.label}>
          Email:
          <input
            className={styles.input}
            type='email'
            name='email'
            value={formData.email}
            onChange={handleChange}
          />
        </label>
        <br />
        <label className={styles.label}>
          Password:
          <input
            className={styles.input}
            type='password'
            name='password'
            value={formData.password}
            onChange={handleChange}
          />
        </label>
        <br />
        <button className={styles.button} type='submit' disabled={mutation.isLoading}>
          {mutation.isLoading ? 'Loading...' : 'Log In'}
        </button>
      </form>
      {(logInError === '') ? null : (<div className={styles.error}>{logInError}</div>)}
    </div>
  );
};

export default LoginForm;
