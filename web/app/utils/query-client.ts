import { QueryCache, QueryClient } from '@tanstack/react-query';

export const queryClient = new QueryClient({
  queryCache: new QueryCache({}),
  defaultOptions: {
    queries: {
      staleTime: 1000 * 60 * 1,
      retry: false,
      // Requests must reach the service worker even when the browser is offline.
      networkMode: 'offlineFirst',
    },
  },
});
