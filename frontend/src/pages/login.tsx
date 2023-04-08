import { useRouter } from "next/router";
import React, { useState } from "react";

import LoginForm from "../modules/login/LoginForm";
import SignUpForm from "../modules/login/SignUpForm";
import withLogin from "../modules/auth/withLogin";

import styles from '../styles/login/LogInPage.module.css'

const LoginPage = () => {
  const router = useRouter();
  const [isNewUser, setIsNewUser] = useState<boolean>(false);

  const handleLoginSuccess = () => {
    router.push("/");
  };

  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Welcome to Points are Bad!</h1>
        { 
          isNewUser ?
          <div>
            <h1 className={styles.h1}>Sign up page</h1>
            <SignUpForm onSuccess={handleLoginSuccess} />
          </div> :
          <div>
            <h1 className={styles.h1}>Login</h1>
            <LoginForm onSuccess={handleLoginSuccess} />
          </div>          
        }
      <div>
        {
          isNewUser ?
          <div className={styles.buttonContainer}>
            <h3 className={styles.h3}>Already got an account?</h3>
            <button className={styles.button} onClick={() => {
              setIsNewUser((prev) => !prev);
            }}>
              Log In
            </button>
          </div> :
          <div className={styles.buttonContainer}>
            <h3 className={styles.h3}>New here?</h3>
            <br />
            <button className={styles.button} onClick={() => {
              setIsNewUser((prev) => !prev);
            }}>
              Create an account
            </button>
          </div>
        }
      </div>
    </div>
  );
};

export default withLogin(LoginPage);