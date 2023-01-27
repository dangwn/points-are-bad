import React from "react";

import headerStyles from "../../styles/header.module.css";

interface userStatProps {
  statKey: string,
  value: string | number | boolean
};

const HeaderUserStat: React.FC<userStatProps> = ({ statKey, value }) => {
  return (
    <div className={headerStyles.headerUserStat}>
      <div className={headerStyles.headerUserStatValue}>{value}</div>
      <div className={headerStyles.headerUserStatKey}>{statKey}</div>
    </div>
  )
};

export default HeaderUserStat;