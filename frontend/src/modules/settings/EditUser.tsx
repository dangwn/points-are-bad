import React, { useState } from 'react';
import { useMutation } from 'react-query';

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
  const [message, setMessage] = useState<string>('');
  const [error, setError] = useState<string>('');
  const [newUsername, setNewUsername] = useState<string>('');
  const [passwordForm, setPasswordForm] = useState<newPasswordForm>({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  });

  const submitUsernameMutation = useMutation(
    async () => {
      try {
        await updateUsername(newUsername);
        setNewUsername('');
        setIsEditingUsername(false);
        setError('');
        setMessage('Username updated successfully!');
      } catch {
        setError('Unable to update username');
      }
  });

  const submitPasswordMutation = useMutation(
    async () => {
      try {
        await updatePassword(passwordForm.oldPassword, passwordForm.newPassword);
        setPasswordForm({
          oldPassword: '',
          newPassword: '',
          confirmPassword: ''
        });
        setIsEditingPassword(false);
        setError('');
        setMessage('Password updated successfully!');
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
          onClick={() => {
            setIsEditingUsername(true);
            setMessage('');
          }}
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
                setMessage('');
              }}
            >
              Cancel
            </button>
        </> :
        <button 
          className={styles.button} 
          onClick={() => {
            setIsEditingPassword(true);
            setMessage('');
          }}
        >
          Change Password
        </button> 
      }
      {
        message && 
        <p className={styles.message}>{message}</p>
      }
      {
        error &&
        <p className={styles.error}>{error}</p>
      }
    </div>
  );
};

export default EditUser;