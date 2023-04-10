import React, { useEffect, useState } from 'react';
import { useMutation } from 'react-query';
import { useRouter } from 'next/router';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/login/SignUpForm.module.css';

type SignUpFormProps = {
  onSuccess: () => void;
};

type VerificationToken = {
  access_token: string,
  token_type: string
}

type SignUpData = {
  username: string,
  password: string,
  confirmPassword: string,
};

const SignUpForm: React.FC<SignUpFormProps> = ({ onSuccess }: SignUpFormProps) => {
  const [email, setEmail] = useState<string>('');
  const [verificationToken, setVerificationToken] = useState<string>('');
  const [formData, setFormData] = useState<SignUpData>({
    username: '',
    password: '',
    confirmPassword: '',
  });
  const [signUpError, setSignUpError] = useState<string>('');
  const router = useRouter();

  useEffect(() => {
    const { token } = router.query;
    if (typeof(token) === 'string') {
      setVerificationToken(token);
    };
  }, [router.query]);

  const createVerificationToken = useMutation(
    async () => {
      const response = await fetch(
        `${API_HOST}/auth/verify/?email=${email}`,
        {
          method: 'POST',
          headers: {
            'accept': 'application/json'
          }
        }
      );
      if (!response.ok) {
        if (response.status === 400) {
          const { detail } = await response.json()
          setSignUpError(detail);
        }
        throw new Error('Could not sign up');
      };

      // @TODO: Make this display "check email" message
      const data: VerificationToken = await response.json();
      router.push(`/signup?token=${data.access_token}`)
    }
  )

  const createUser = useMutation(
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
            body: JSON.stringify({
              token: verificationToken,
              username: formData.username,
              password: formData.password
            }),
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

  const handleCreateVerificationTokenSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSignUpError('');
    createVerificationToken.mutate();
  };

  const handleCreateUserSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSignUpError('');
    createUser.mutate();
  };

  const handleCreateVerificationTokenChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value)
  };

  const handleCreateUserChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData((prevData) => ({
      ...prevData,
      [e.target.name]: e.target.value,
    }));
  };

  return (
    <div className={styles.container}>
      <div>
        {
          (verificationToken === '') ?
          <form onSubmit={handleCreateVerificationTokenSubmit}>
            <label className={styles.label}>
              Email:
              <input 
                className={styles.input}
                type='email'
                name='email'
                value={email}
                onChange={handleCreateVerificationTokenChange}
              />
            <button 
              className={styles.button} 
              type='submit' 
              disabled={createVerificationToken.isLoading}
            >
              {createVerificationToken.isLoading ? 'Loading...': 'Verify Email'}
            </button>
            </label>
          </form> :
          <form onSubmit={handleCreateUserSubmit}>
            <label className={styles.label}>
              Username:
              <input 
                className={styles.input}
                type='text'
                name='username'
                value={formData.username}
                onChange={handleCreateUserChange}
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
                onChange={handleCreateUserChange}
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
                onChange={handleCreateUserChange}
              />
            </label>
            <button 
              className={styles.button} 
              type='submit' 
              disabled={createUser.isLoading}
            >
              {createUser.isLoading ? 'Loading...': 'Sign Up'}
            </button>
          </form>
        }
      </div>
      {(signUpError === '') ? null : (<div className={styles.error}>{signUpError}</div>)}
    </div>
  )
}

export default SignUpForm;