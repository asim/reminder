import {
  Links,
  Meta,
  Outlet,
  Scripts,
  ScrollRestoration
} from 'react-router';

import { QueryClientProvider } from '@tanstack/react-query';
import { Loader2Icon } from 'lucide-react';
import { NavigationProgress } from '~/components/interface/navigation-progress';
import { PageError } from '~/components/interface/page-error';
import { queryClient } from '~/utils/query-client';
import type { Route } from './+types/root';
import './app.css';

export const links: Route.LinksFunction = () => [
  {
    rel: 'preload',
    href: '/fonts/arabic.otf',
    as: 'font',
    type: 'font/otf',
    crossOrigin: 'anonymous',
  },
  {
    rel: 'manifest',
    href: '/manifest.webmanifest',
  },
];

export function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html lang='en'>
      <head>
        <meta charSet='utf-8' />
        <meta name='viewport' content='width=device-width, initial-scale=1, interactive-widget=resizes-content' />
        <meta name='theme-color' content='#ffffff' />
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
  return <PageError error={error} />;
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
