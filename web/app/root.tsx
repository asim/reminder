import {
  isRouteErrorResponse,
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration,
} from 'react-router';

import type { Route } from './+types/root';
import './app.css';
import { QueryClientProvider } from '@tanstack/react-query';
import { queryClient } from '~/utils/query-client';
import { Loader2Icon } from 'lucide-react';
import { NavigationProgress } from '~/components/interface/navigation-progress';

export const links: Route.LinksFunction = () => [
  {
    rel: 'preload',
    href: '/fonts/arabic.otf',
    as: 'font',
    type: 'font/otf',
    crossOrigin: 'anonymous',
  },
];

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang='en'>
      <head>
        <meta charSet='utf-8' />
        <meta name='viewport' content='width=device-width, initial-scale=1' />
        <Meta />
        <Links />
      </head>
      <body>
        <QueryClientProvider client={queryClient}>
          {children}
          <ScrollRestoration />
          <Scripts />
          <NavigationProgress />
        </QueryClientProvider>
      </body>
    </html>
  );
}

export default function App() {
  return <Outlet />;
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  let message = 'Oops!';
  let details = 'An unexpected error occurred.';
  let stack: string | undefined;

  if (isRouteErrorResponse(error)) {
    message = error.status === 404 ? '404' : 'Error';
    details =
      error.status === 404
        ? 'The requested page could not be found.'
        : error.statusText || details;
  } else if (import.meta.env.DEV && error && error instanceof Error) {
    details = error.message;
    stack = error.stack;
  }

  return (
    <main className='pt-16 p-4 container mx-auto'>
      <h1>{message}</h1>
      <p>{details}</p>
      {stack && (
        <pre className='w-full p-4 overflow-x-auto'>
          <code>{stack}</code>
        </pre>
      )}
    </main>
  );
}

export function HydrateFallback() {
  return (
    <div className='bg-opacity-75 fixed top-0 left-0 z-[100] flex h-full w-full items-center justify-center bg-white'>
      <div className='flex items-center justify-center rounded-lg border border-gray-200 bg-white px-4 py-2'>
        <Loader2Icon className='size-4 animate-spin stroke-[2.5] text-gray-500' />
        <span className='ml-2 text-sm text-black'>
          please wait&nbsp;
          <span className='animate-pulse'>...</span>
        </span>
      </div>
    </div>
  );
}
