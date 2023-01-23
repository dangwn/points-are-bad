import { Html, Head, Main, NextScript } from 'next/document'
import React from 'react'

interface docProps {};

const Document: React.FC<docProps> = () => {
  return (
    <Html lang="en">
      <Head />
      <body>
        <Main />
        <NextScript />
      </body>
    </Html>
  );
};

export default Document;