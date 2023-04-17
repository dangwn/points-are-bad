import React from "react";
import { useRouter } from "next/router";

import SignUpForm from "./SignUpForm";
import styles from '../../styles/login/SignUpForm.module.css'

interface SignUpPageProps {}

const SignUpPage: React.FC<SignUpPageProps> = ({}) => {
  const router = useRouter();

  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Welcome to Points are Bad!</h1>
      <div>
        <h1 className={styles.h1}>Sign Up</h1>
        <SignUpForm onSuccess={() => {
          router.push('/');
        }} />
      </div>        
      <div>
        <div className={styles.buttonContainer}>
          <h3 className={styles.h3}>Already got an account?</h3>
          <button className={styles.button} onClick={() => {
            router.push('/login');
          }}>
            Log In
          </button>
        </div>
      </div>
    </div>
  );
};

export default SignUpPage;