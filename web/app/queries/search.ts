import { httpGet, httpPost } from '~/utils/http';
import { queryClient } from '~/utils/query-client';

type SearchReference = {
  text: string;
  score: number;
  metadata: Record<string, string>;
};

export type SearchSource = 'all' | 'quran' | 'hadith' | 'tafsir' | 'names';
export type SearchResponseType = {
  source?: string;
  mode?: 'search' | 'ask';
  q: string;
  answer: string;
  references: SearchReference[];
};

export type SearchHistoryType = {
  session: string;
  history: string[];
};

// Fast search: fetches references only, no LLM summarisation.
export const searchOptions = (query: string, source: SearchSource = 'all') => ({
  queryKey: ['search', query, source],
  queryFn: async () => {
    if (!query.trim()) {
      return null;
    }

    const data = await httpPost<SearchResponseType>('/api/search', {
      q: query,
      source,
      summarise: false,
    });

    return data;
  },
  enabled: !!query.trim(),
});

// Slow search: triggers LLM summarisation for the given query.
export const searchSummaryOptions = (query: string, source: SearchSource = 'all') => ({
  queryKey: ['search-summary', query, source],
  staleTime: Infinity,
  retry: false,
  refetchOnWindowFocus: false,
  refetchOnReconnect: false,
  refetchOnMount: false,
  queryFn: async () => {
    if (!query.trim()) {
      return null;
    }

    const data = await httpPost<SearchResponseType>('/api/search', {
      q: query,
      source,
      summarise: true,
    });

    queryClient.invalidateQueries({
      queryKey: getSearchHistoryOptions().queryKey,
    });

    return data;
  },
  enabled: !!query.trim(),
});

export const getSearchHistoryOptions = () => ({
  queryKey: ['search-history'],
  queryFn: async () => {
    const data = await httpGet<SearchHistoryType>('/api/search');
    return data;
  },
});
