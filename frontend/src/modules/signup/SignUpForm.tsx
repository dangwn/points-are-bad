import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { useMutation } from 'react-query';

import Loading from '../shared/Loading';
import { setAccessToken } from '../../lib/accessToken';
import { createUser, getVerificationToken } from '../../lib/requests';
import styles from '../../styles/login/SignUpForm.module.css';
import type { SignUpData } from '../../types/auth';
import type { Token } from '../../types/token';

interface SignUpFormProps {
  onSuccess: () => void
}

const SignUpForm: React.FC<SignUpFormProps> = ({ onSuccess }) => {
  const [email, setEmail] = useState<string>('');
  const [verificationToken, setVerificationToken] = useState<string>('');
  const [formData, setFormData] = useState<SignUpData>({
    username:'',
    password:'',
    confirmPassword:''
  })
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
      const token: Token = await getVerificationToken(email);

      // @TODO: Make this display "check email" message
      router.push(`/signup?token=${token.access_token}`);
    }
  )

  const signUpUser = useMutation(
    async () => {
      if (formData.password !== formData.confirmPassword) {
        setSignUpError('Passwords do not match!');
      } else if (formData.password.length < 8) {
        setSignUpError('Password must be at least 8 characters!');
      } else if (formData.username.length < 3) {
        setSignUpError('Username must be at least 3 characters long!');
      } else {
        try {
          const data: Token = await createUser(
            verificationToken,
            formData.username,
            formData.password
          );
          setAccessToken(data.access_token);
          return;
        } catch {
          setSignUpError('Could not sign user up. Please start again.')
        };
      };
      throw new Error('Could not sign up user.')
    },
    {
      onSuccess
    }
  )

  const handleCreateVerificationTokenSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSignUpError('');
    createVerificationToken.mutate();
  };

  const handleCreateUserSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSignUpError('');
    signUpUser.mutate();
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
              {createVerificationToken.isLoading ? 'Loading...' : 'Verify Email'}
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
              disabled={signUpUser.isLoading}
            >
              {signUpUser.isLoading ? 'Loading...': 'Sign Up'}
            </button>
          </form>
        }
      </div>
      {(signUpError === '') ? null : (<div className={styles.error}>{signUpError}</div>)}
    </div>
  )
}

export default SignUpForm;