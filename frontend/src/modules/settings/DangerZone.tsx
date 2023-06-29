import { useRouter } from 'next/router';
import React, { useState } from 'react';
import { useMutation, useQueryClient } from 'react-query';

import { deleteCurrentUser } from '@/lib/requests';
import styles from '@/styles/settings/DangerZone.module.css';

const DangerZone: React.FC = () =>{
  const queryClient = useQueryClient();
  const router = useRouter();
  const [ showConfirmation, setShowConfirmation ] = useState<boolean>(false)
  const [ deleteUserError, setDeleteUserError ] = useState<string>('')

  const handleDeleteUser = useMutation(async () => {
    const response = await deleteCurrentUser();
    if (!response.ok) {
      if (response.status === 403) {
        setDeleteUserError('You are the only remaining admin user. Please make someone else admin before deleting your account.');
      } else {
        setDeleteUserError('Unable to delete user.');
      }

      return
    }

    queryClient.removeQueries();

    router.push('/login');
  })

  return (
    <>
    <div className={styles.dangerZoneContainer}>
      <h2 className={styles.h2}>Danger Zone</h2>
      <button className={styles.button} onClick={() => setShowConfirmation(true)}>Delete Account</button>
    </div>
    {
      showConfirmation ?
      <div className={styles.deleteUserPopup}>
        <div className={styles.deleteUserPopupContainer}>
          <h2 className={styles.h2}>Delete account? This action cannot be undone.</h2>
          <button className={styles.popupButton} onClick={() => handleDeleteUser.mutate()}>Confirm</button>
          <br />
          <button className={styles.popupButton} onClick={() => {
            setShowConfirmation(false);
            setDeleteUserError('');
          }}>Cancel</button>
          {
            deleteUserError ?
            <p className={styles.deleteUserError}>{deleteUserError}</p> :
            null
          }
        </div>
      </div> :
      null
    }
    </>
  )
}

export default DangerZone;