import { useEffect, useState } from "react";
import Cookie from "js-cookie";

import styles from '../../styles/cookie/CookiePopup.module.css'

const CookiePopup = () => {
  const [showPopup, setShowPopup] = useState(false);

  useEffect(() => {
    const acceptCookie = Cookie.get("PAB-Accept-Cookie");

    if (!acceptCookie) {
      setShowPopup(true);
    }
  }, []);

  const handleAcceptCookies = () => {
    Cookie.set("PAB-Accept-Cookie", "true", { expires: 365 });
    setShowPopup(false);
  };

  return showPopup ? (
    <div className={styles.cookiePopup}>
      <div className={styles.cookiePopupContainer}>
        <p className={styles.p}>
          Points are Bad only uses necessary cookies.
        </p>
        <p className={styles.p}>
          Click &quot;Accept&quot; to continue.
        </p>
        <button className={styles.button} onClick={handleAcceptCookies}>Accept</button>
      </div>
    </div>
  ) : null;
};

export default CookiePopup;
