import { AppProps } from 'next/app';
import { QueryClientProvider } from 'react-query';
import { queryClient } from '../lib/queryClient';
import CookiePopup from '../modules/cookie/CookiePopup';

function MyApp({ Component, pageProps }: AppProps) {
  return (
    <>
      <QueryClientProvider client={queryClient}>
        <Component {...pageProps} />
      </QueryClientProvider>
      <CookiePopup />
    </>
  );
}

export default MyApp;