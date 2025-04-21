import { isRouteErrorResponse } from 'react-router';

import { Frown } from 'lucide-react';
import { FetchError } from '~/utils/http';

type PageErrorProps = {
  error: unknown;
};

export function PageError({ error }: PageErrorProps) {
  let title = 'Oops!';
  let subtitle = 'An unexpected error occurred.';

  if (error instanceof FetchError) {
    if (error.status === 404) {
      title = 'Page not found';
      subtitle = 'Please make sure the page you are looking for exists.';
    } else {
      title = error.message;
      subtitle = 'Please consider refreshing the page or trying again later.';
    }
  } else if (isRouteErrorResponse(error)) {
    title = error.statusText;
    subtitle = error.data;
  } else if (error instanceof Error) {
    title = error.message;
    subtitle = 'Please consider refreshing the page or trying again later.';
  }

  return (
    <div className='max-w-4xl mx-auto flex flex-col items-center justify-center h-screen'>
      <div className='flex flex-col items-center justify-center'>
        <Frown className='size-18 mb-6 text-gray-300' />
        <h1 className='text-3xl mb-2 font-bold'>{title}</h1>
        <p className='text-lg text-gray-500'>{subtitle}</p>
      </div>
    </div>
  );
}
