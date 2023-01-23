import React from "react";
import { signIn } from "next-auth/react";

import { CLIENT_HOST } from "@/lib/constants";

interface ProviderFormProps {};

const ProviderForm: React.FC<ProviderFormProps> = () => {
  const signInCallbackUrl = CLIENT_HOST;
  
  return (
    <>
      <div className='login-signup-form'>
        <div>
          <h3>Welcome to Points are Bad</h3>
          <hr />
          <div className='auth-button'>
            <button onClick={() => signIn(
              'github',
              {
                callbackUrl:signInCallbackUrl
              }
            )}>Continue with GitHub</button>
          </div>
          <div className='auth-button'>
            <button onClick={() => signIn(
              'google',
              {
                callbackUrl:signInCallbackUrl
              }
            )}>Continue with Google</button>
          </div>
        </div>
      </div>
    </>
  );
};

export default ProviderForm;