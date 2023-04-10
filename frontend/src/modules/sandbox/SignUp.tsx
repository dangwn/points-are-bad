import { useEffect, useState } from 'react';
import { useMutation } from 'react-query';
import { useRouter } from 'next/router';

import { API_HOST } from '@/lib/constants';

interface SignUpFormData {
  email: string;
  verificationCode: string|URL;
  username: string;
  password: string;
  matchPassword: string;
}

const sendVerificationCode = async (email: string) => {
  const url = new URL(`http://localhost:3000/sandbox`)
  url.searchParams.set('email', 'dan@dan.com')
  url.searchParams.set('verificationCode','123456')
  return url
}

const createUser = async ({
  email,
  verificationCode,
  username,
  password,
  matchPassword,
}: SignUpFormData) => {

}

export default function SignUp() {
  const [email, setEmail] = useState('');
  const [verificationCode, setVerificationCode] = useState<string|URL>('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [matchPassword, setMatchPassword] = useState('');
  const router = useRouter();

  const sendCodeMutation = useMutation(sendVerificationCode);
  const createUserMutation = useMutation(createUser);

  useEffect(() => {
    const {email, verificationCode} = router.query;
    setVerificationCode(`${email}||${verificationCode}`)
  }, [router.query])

  const handleSendCode = async () => {
    try {
      const link = await sendCodeMutation.mutateAsync(email);
      // email sent successfully, display message to user
      router.push(link)
    } catch (error) {
      // handle error
    }
  };

  const handleCreateUser = async () => {
    try {
      await createUserMutation.mutateAsync({
        email,
        verificationCode,
        username,
        password,
        matchPassword,
      });
      router.push('/login'); // redirect user to login page
    } catch (error) {
      // handle error
    }
  };

  const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    handleCreateUser();
  };

  return (
    <div>
    <form onSubmit={handleSubmit}>
      <label htmlFor="email">Email:</label>
      <input
        id="email"
        type="email"
        value={email}
        onChange={(event) => setEmail(event.target.value)}
        required
      />
      <button type="button" onClick={handleSendCode}>
        Send Verification Code
      </button>
      {/* {verificationCode && (
        <>
          <label htmlFor="username">Username:</label>
          <input
            id="username"
            type="text"
            value={username}
            onChange={(event) => setUsername(event.target.value)}
            required
          />
          <label htmlFor="password">Password:</label>
          <input
            id="password"
            type="password"
            value={password}
            onChange={(event) => setPassword(event.target.value)}
            required
          />
          <label htmlFor="matchPassword">Confirm Password:</label>
          <input
            id="matchPassword"
            type="password"
            value={matchPassword}
            onChange={(event) => setMatchPassword(event.target.value)}
            required
          />
          <button type="submit">Create Account</button>
        </>
      )} */}
    </form>
    <div>{verificationCode.toString()}</div>
    </div>
  );
}