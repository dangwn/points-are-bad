import React, { useEffect, useState } from 'react';
import { useRouter } from 'next/router';
import { useMutation } from 'react-query';

import { setAccessToken } from '@/lib/accessToken';
import { createUser, getVerificationToken } from '@/lib/requests';

import styles from '@/styles/SignUpPage.module.css';
import type { SignUpData } from '@/types/auth';
import type { Token } from '@/types/token';

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
      <div>
        {
          (verificationToken === '') ?
          <div className={styles.signUpForm}>
          <form className={styles.signUpForm} onSubmit={handleCreateVerificationTokenSubmit}>
            <input 
              className={styles.input}
              type='email'
              name='email'
              value={email}
              onChange={handleCreateVerificationTokenChange}
              placeholder='Email'
            />
            <button 
              className={styles.button} 
              type='submit' 
              disabled={createVerificationToken.isLoading}
            >
              {createVerificationToken.isLoading ? 'Loading...' : 'Verify Email'}
            </button>
          </form>
          <hr className={styles.signUpHR}/>
          <p className={styles.h3}>Already got an account?</p>
          <button className={styles.button} onClick={() => {
            router.push('/login');
          }}>
            Log In
          </button>
          </div> :
          <form className={styles.signUpForm} onSubmit={handleCreateUserSubmit}>
              <input 
                className={styles.input}
                type='text'
                name='username'
                value={formData.username}
                onChange={handleCreateUserChange}
                placeholder='Username'
              />
              <input 
                className={styles.input}
                type='password'
                name='password'
                value={formData.password}
                onChange={handleCreateUserChange}
                placeholder='Password'
              />
              <input 
                className={styles.input}
                type='password'
                name='confirmPassword'
                value={formData.confirmPassword}
                onChange={handleCreateUserChange}
                placeholder='Confirm Password'
              />
            <button 
              className={styles.button} 
              type='submit' 
              disabled={signUpUser.isLoading}
            >
              {signUpUser.isLoading ? 'Loading...': 'Create Account'}
            </button>
          </form>
        }
        {(signUpError === '') ? null : (<div className={styles.error}>{signUpError}</div>)}
      </div>
  )
}

export default SignUpForm;