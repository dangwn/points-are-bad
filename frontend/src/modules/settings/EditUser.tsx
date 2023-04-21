import React, { useState } from 'react';
import { useMutation } from 'react-query';

import { updateUsername } from '@/lib/requests';
import styles from '@/styles/settings/EditUser.module.css';

interface EditUserProps {
  username: string
}

const EditUser: React.FC<EditUserProps> = ({ username }) => {
  const [isEditingUsername, setIsEditingUsername] = useState<boolean>(false);
  const [newUsername, setNewUsername] = useState<string>('');

  const submitUsernameMutation = useMutation(
    async () => {
      await updateUsername(newUsername);
  });

  const handleUsernameSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setNewUsername('');
    setIsEditingUsername(false);
    submitUsernameMutation.mutate();
  }

  return (
    <div className={styles.editUserContainer}>
      <h2 className={styles.h2}>Edit Profile</h2>
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
              onClick={() => setIsEditingUsername(false)}
            >
              Cancel
            </button>
        </> :
        <button className={styles.button} onClick={() => setIsEditingUsername(true)}>Edit Username</button> 
      }
    </div>
  );
};

export default EditUser;