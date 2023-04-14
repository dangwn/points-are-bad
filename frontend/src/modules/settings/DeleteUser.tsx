import React, { useState } from 'react';
import { useMutation, useQueryClient} from 'react-query';
import { useRouter } from 'next/router';

import { API_HOST } from '../../lib/constants';
import styles from '../../styles/settings/DeleteUser.module.css'

const DeleteUser: React.FC = () =>{
  const queryClient = useQueryClient();
  const router = useRouter();
  const [ showConfirmation, setShowConfirmation ] = useState<boolean>(false)

  const handleDeleteUser = useMutation(async () => {
    const accessToken = localStorage.getItem('access_token');
    await fetch(`${API_HOST}/user/`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${accessToken}`
      },
      credentials: 'include',
    });

    queryClient.removeQueries();

    router.push('/login');
  })

  return (
    <div>
      <button className={styles.button} onClick={() => setShowConfirmation(true)}>Delete Account</button>
      {
        showConfirmation ?
        <div className={styles.modal}>
          <div className={styles.modalContent}>
            <h2>Are you sure you want to delete your account?</h2>
            <button className={styles.button} onClick={() => handleDeleteUser.mutate()}>Confirm</button>
            <button className={styles.button} onClick={() => setShowConfirmation(false)}>Cancel</button>
          </div>
        </div> :
        null
      }
    </div>
  )
}

export default DeleteUser;