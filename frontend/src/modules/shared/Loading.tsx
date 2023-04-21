import React from "react";

import styles from '@/styles/shared/Loading.module.css'

const Loading: React.FC = () => {
  return (
    <div className={styles.loadingContainer}>
      <div className={styles.spinner}>
        <div className={styles.spinnerDot}></div>
        <div className={styles.spinnerDot}></div>
        <div className={styles.spinnerDot}></div>
        <div className={styles.spinnerDot}></div>
      </div>
    </div>
  );
}

export default Loading;