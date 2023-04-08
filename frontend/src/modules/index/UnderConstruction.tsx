import TestButton from './TestButton'

import styles from '../../styles/index/UnderConstruction.module.css'

const UnderConstruction = () => {
  return (
    <div className={styles.container}>
      <h1 className={styles.h1}>Under construction</h1>
      <TestButton />
    </div>
  );
}

export default UnderConstruction;