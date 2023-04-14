import { useRouter } from "next/router";
import React from "react";

import LoginForm from "../modules/login/LoginForm";
import withLogin from "../modules/auth/withLogin";

import styles from '../styles/login/LogInForm.module.css'

const LoginPage = () => {
  const router = useRouter();

  const handleLoginSuccess = () => {
    router.push("/");
  };

  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Welcome to Points are Bad!</h1>
      <div>
        <h1 className={styles.h1}>Login</h1>
        <LoginForm onSuccess={handleLoginSuccess} />
      </div>          
      <div>
        <div className={styles.buttonContainer}>
          <h3 className={styles.h3}>New here?</h3>
          <button className={styles.button} onClick={() => {
            router.push('/signup');
          }}>
            Create an account
          </button>
        </div>
      </div>
    </div>
  );
};

export default withLogin(LoginPage);