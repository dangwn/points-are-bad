import React, { useState } from 'react';
import { useMutation } from 'react-query';
import { useRouter } from 'next/router';

import { updateUsername, updatePassword } from '@/lib/requests';
import styles from '@/styles/settings/EditUser.module.css';

interface EditUserProps {
  username: string
}

type newPasswordForm = {
  oldPassword: string,
  newPassword: string,
  confirmPassword: string
}

const EditUser: React.FC<EditUserProps> = ({ username }) => {
  const [isEditingUsername, setIsEditingUsername] = useState<boolean>(false);
  const [isEditingPassword, setIsEditingPassword] = useState<boolean>(false);
  const [error, setError] = useState<string>('');
  const [newUsername, setNewUsername] = useState<string>('');
  const [passwordForm, setPasswordForm] = useState<newPasswordForm>({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  });
  const router = useRouter();

  const submitUsernameMutation = useMutation(
    async () => {
      try {
        await updateUsername(newUsername);
        router.reload();
      } catch {
        setError('Unable to update username');
      }
  });

  const submitPasswordMutation = useMutation(
    async () => {
      try {
        await updatePassword(passwordForm.oldPassword, passwordForm.newPassword);
        router.reload();
      } catch (e: any) {
        setError(e.message);
      }
    }
  )

  const handlePasswordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setPasswordForm((prevData) => ({
      ...prevData,
      [e.target.name]: e.target.value,
    }));
  };

  const handleUsernameSubmit = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    if (newUsername.length < 3){
      setError('Username must be more than three characters!')
    } else {
      submitUsernameMutation.mutate();
    }
  }

  const handlePasswordSubmit = (e: React.FormEvent<HTMLFormElement>): void => {
    e.preventDefault();
    if (passwordForm.newPassword != passwordForm.confirmPassword) {
      setError('New passwords do not match!');
    } else if (passwordForm.newPassword.length < 8) {
      setError('New password must be at least 8 characters!');
    } else {
      submitPasswordMutation.mutate();
    }
  }

  return (
    <div className={styles.editUserContainer}>
      <h2 className={styles.h2}>Profile Settings</h2>
      <h3>Username</h3>
      {
        isEditingUsername ?
        <>
          <form 
            className={styles.form} 
            onSubmit={handleUsernameSubmit}
          >
            <input 
              className={styles.input}
              type='text'
              placeholder={username}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                setNewUsername(e.target.value);
              }}
            />
            <button 
              className={styles.button}
              type='submit'
            >
              Save
            </button>
            </form>
            <button 
              className={styles.button} 
              onClick={() => {
                setIsEditingUsername(false);
                setError('');
              }}
            >
              Cancel
            </button>
        </> :
        <button 
          className={styles.button} 
          onClick={() => setIsEditingUsername(true)}
        >
          Change Username
        </button> 
      }
      <h3>Password</h3>
      {
        isEditingPassword ?
        <>
          <form 
            className={styles.form} 
            onSubmit={handlePasswordSubmit}
          >
            <input 
              className={styles.input}
              type='password'
              placeholder='Old Password'
              name='oldPassword'
              onChange={handlePasswordChange}
            />
            <input 
              className={styles.input}
              type='password'
              placeholder='New Password'
              name='newPassword'
              onChange={handlePasswordChange}
            />
            <input 
              className={styles.input}
              type='password'
              placeholder='Confirm Password'
              name='confirmPassword'
              onChange={handlePasswordChange}
            />
            <button 
              className={styles.button}
              type='submit'
            >
              Save
            </button>
            </form>
            <button 
              className={styles.button} 
              onClick={() => {
                setIsEditingPassword(false);
                setError('');
              }}
            >
              Cancel
            </button>
        </> :
        <button 
          className={styles.button} 
          onClick={() => setIsEditingPassword(true)}
        >
          Change Password
        </button> 
      }
      {
        error &&
        <p className={styles.error}>{error}</p>
      }
    </div>
  );
};

export default EditUser;