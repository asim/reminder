import { useQuery } from '@tanstack/react-query';

// One shared request across search and every reader action. Unknown availability
// is conservative; a failed refresh retains the last successful configuration.
export function useModelAvailable() {
  const { data } = useQuery({
    queryKey: ['capabilities'],
    queryFn: async () => {
      const response = await fetch('/api/capabilities', { cache: 'no-store' });
      if (!response.ok) throw new Error('Capabilities unavailable');
      return await response.json() as { model: boolean };
    },
    staleTime: 60_000,
    retry: false,
    refetchOnWindowFocus: false,
  });
  return data?.model === true;
}
