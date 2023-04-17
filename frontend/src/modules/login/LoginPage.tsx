import React from "react";
import { useRouter } from "next/router";

import LoginForm from "./LoginForm";
import styles from '../../styles/login/LogInForm.module.css'

interface LoginPageProps {}

const LoginPage: React.FC<LoginPageProps> = ({}) => {
  const router = useRouter();

  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Welcome to Points are Bad!</h1>
      <div>
        <h1 className={styles.h1}>Login</h1>
        <LoginForm onSuccess={() => router.push('/')} />
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

export default LoginPage;