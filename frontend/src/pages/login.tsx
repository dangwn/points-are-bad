import React, { useEffect } from "react";
import { useRouter } from "next/router";
import { useSession } from "next-auth/react";
import ProviderForm from "@/modules/login/ProviderForm";

interface loginProps {};

const LoginPage: React.FC<loginProps> = ({}) => {
  const { data: session, status} = useSession();
  const { push } = useRouter();

  useEffect(() => {
    if (session) {
      push({
        pathname: '/'
      });
    };

    return (): void => {};
  }, [session, push]);

  if (status === 'unauthenticated') {
    return (
      <>
        <ProviderForm />
      </>
    );
  }

  return <p>loading...</p>
};

export default LoginPage;