import { useState } from 'react';
import { useMutation } from 'react-query';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/login/LogInForm.module.css'

type LoginFormProps = {
  onSuccess: () => void;
};

type LoginData = {
  email: string;
  password: string;
};

const LoginForm = ({ onSuccess }: LoginFormProps) => {
  const [formData, setFormData] = useState<LoginData>({
    email: '',
    password: '',
  });
  const [logInError, setLogInError] = useState<string>('');

  const mutation = useMutation(
    async () => {
      const response = await fetch(
        `${API_HOST}/auth/login/`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'accept': 'application/json'
          },
          body: JSON.stringify(formData),
          credentials: 'include'
        }
      );

      if (!response.ok) {
        if (response.status === 401) {
          setLogInError('Incorrect email or password');
        }
        throw new Error('Invalid credentials');
      };

      const data = await response.json();
      localStorage.setItem('access_token', data.access_token);
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
