import React, { useEffect } from 'react';
import { useRouter } from 'next/router';

import withAuth from '../modules/auth/withAuth';

const ErrorPage: React.FC = () => {
  const router = useRouter();

  useEffect(() => {
    router.push('/')
  }, [router]);
  
  return <></>
}

export default withAuth(ErrorPage);