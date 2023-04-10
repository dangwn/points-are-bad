import { useRouter } from "next/router";
import React from "react";

import SignUpForm from "../modules/login/SignUpForm";
import withLogin from "../modules/auth/withLogin";

import styles from '../styles/login/LogInPage.module.css'

const LoginPage = () => {
  const router = useRouter();

  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Welcome to Points are Bad!</h1>
      <div>
        <h1 className={styles.h1}>Sign up page</h1>
        <SignUpForm onSuccess={() => {
          router.push('/')
        }} />
      </div>        
      <div>
      <div className={styles.buttonContainer}>
        <h3 className={styles.h3}>Already got an account?</h3>
        <button className={styles.button} onClick={() => {
          router.push('/login')
        }}>
          Log In
        </button>
      </div>

      </div>
    </div>
  );
};

export default withLogin(LoginPage);