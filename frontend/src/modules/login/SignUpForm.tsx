import React, { useState } from 'react';
import { useMutation } from 'react-query';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/login/SignUpForm.module.css'

type SignUpFormProps = {
  onSuccess: () => void;
};

type SignUpData = {
  email: string;
  username: string;
  password: string;
  confirmPassword: string;
};

const SignUpForm: React.FC<SignUpFormProps> = ({ onSuccess }: SignUpFormProps) => {
  const [formData, setFormData] = useState<SignUpData>({
    email: '',
    username: '',
    password: '',
    confirmPassword: '',
  });
  const [signUpError, setSignUpError] = useState<string>('');

  const mutation = useMutation(
    async () => {
      if (formData.password !== formData.confirmPassword) {
        setSignUpError('Passwords do not match!');
        throw new Error('Passwords do not match');
      } else if (formData.password.length < 8) {
        setSignUpError('Password must be at least 8 characters!');
        throw new Error('Password not long enough')
      } else {
        const response = await fetch(
          `${API_HOST}/user/`,
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
          if (response.status === 403){
            const { detail } = await response.json()
            setSignUpError(detail);
          }
          throw new Error('Could not sign up');
        };

        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
      }
    },
    {
      onSuccess,
    }
  );

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSignUpError('');
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
          Username:
          <input 
            className={styles.input}
            type='text'
            name='username'
            value={formData.username}
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
        <label className={styles.label}>
          Confirm Password:
          <input 
            className={styles.input}
            type='password'
            name='confirmPassword'
            value={formData.confirmPassword}
            onChange={handleChange}
          />
        </label>
        <br />
        <button className={styles.button} type='submit' disabled={mutation.isLoading}>
          {mutation.isLoading ? 'Loading...': 'Sign Up'}
        </button>
      </form>
      {(signUpError === '') ? null : (<div className={styles.error}>{signUpError}</div>)}
    </div>
  )
}

export default SignUpForm;