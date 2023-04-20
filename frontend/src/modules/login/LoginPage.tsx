import React from "react";
import { useRouter } from "next/router";

import LoginForm from "./LoginForm";
import styles from '@/styles/LoginPage.module.css'

interface LoginPageProps {}

const LoginPage: React.FC<LoginPageProps> = ({}) => {
  const router = useRouter();

  return (
    <div className={styles.container}>
    <div className={styles.loginPageContainer}>
      <h1 className={styles.loginHeader}>Points Are Bad</h1>
      <LoginForm onSuccess={() => router.push('/')} />
      <hr className={styles.loginHR} />
      <div className={styles.loginButtonContainer}>
        <button className={styles.button} onClick={() => {
          router.push('/signup');
        }}>
          Create an Account
        </button>
      </div>
    </div></div>
  );
};

export default LoginPage;