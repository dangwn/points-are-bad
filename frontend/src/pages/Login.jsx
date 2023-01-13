import { useState } from 'react';
import { SignUpForm } from '../components/SignUpForm';
import { LogInForm } from '../components/LogInForm';

export const Login = () => {
  const[landingState, setLandingState] = useState('Login');

  let landingForm;
  if (landingState === 'Login') {
    landingForm = <LogInForm />;
  } else {
    landingForm = <SignUpForm />;
  };

  return (
    <>
      <div className='login-signup-form'>
        <div className='login-landing-form'>
          {landingForm}
          <hr />
          <div>
            {(landingState === 'Login') ? "New here?" : null}
          </div>
          <button onClick={() => setLandingState(
              (landingState === 'Login') ? 'Signup' : 'Login'
            )}>
              {(landingState === 'Login') ? 'Create New account' : 'Login'}
          </button>
        </div>
      </div>
    </>
  );
};

