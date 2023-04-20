import React from "react";
import { useRouter } from "next/router";

import SignUpForm from "./SignUpForm";

import styles from '@/styles/SignUpPage.module.css';

interface SignUpPageProps {}

const SignUpPage: React.FC<SignUpPageProps> = ({}) => {
  const router = useRouter();

  return (
    <div className={styles.container}>
      <div className={styles.signUpPageContainer}>
        <h1 className={styles.signUpHeader}>Sign Up</h1>
        <SignUpForm onSuccess={() => {
          router.push('/');
        }} />
      </div>
    </div>
  );
};

export default SignUpPage;