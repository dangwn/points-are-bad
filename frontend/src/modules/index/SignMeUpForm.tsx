import Image from "next/image";
import { signOut, signIn} from "next-auth/react"

import { createNewUser } from '@/lib/requests/auth';

import type { Session } from "next-auth"

interface signMeUpFormProps {
  session: Session
}

const SignMeUpForm: React.FC<signMeUpFormProps> = ({ session }) => {
  const handleNewUserRequest = async (): Promise<void> => {
    if (session) {
      const newUserData = await createNewUser(
        session.user.accessToken,
        session.user.provider,
        {
          display_name: session.user.name,
          avatar: session.user.image
        }
      );

      if (newUserData.displayName === ''){
        signOut();
      } else {
        signIn(session.user.provider);
      };
    } else {
      signOut();
    };
  };

  return (
    <>
      <h3>No user linked to your {session.user.provider} account. Would you like to sign up?</h3>
      <div>{`Display name: ${session.user.name}`}</div>
      <Image src={session.user.image} alt="New user's avatar" width={300} height={300}/>
      <button onClick={handleNewUserRequest}>Sign me up!</button>
      <button onClick={() => signOut()}>I don't want to have fun</button>
    </>
  );
};


export default SignMeUpForm