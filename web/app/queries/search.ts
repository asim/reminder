import { httpGet, httpPost } from '~/utils/http';
import { queryClient } from '~/utils/query-client';

type SearchReference = {
  text: string;
  score: number;
  metadata: Record<string, string>;
};

export type SearchResponseType = {
  q: string;
  answer: string;
  references: SearchReference[];
};

export type SearchHistoryType = {
  session: string;
  history: string[];
};

// Fast search: fetches references only, no LLM summarisation.
export const searchOptions = (query: string) => ({
  queryKey: ['search', query],
  queryFn: async () => {
    if (!query.trim()) {
      return null;
    }

    const data = await httpPost<SearchResponseType>('/api/search', {
      q: query,
      summarise: false,
    });

    return data;
  },
  enabled: !!query.trim(),
});

// Slow search: triggers LLM summarisation for the given query.
export const searchSummaryOptions = (query: string) => ({
  queryKey: ['search-summary', query],
  queryFn: async () => {
    if (!query.trim()) {
      return null;
    }

    const data = await httpPost<SearchResponseType>('/api/search', {
      q: query,
      summarise: true,
    });

    queryClient.invalidateQueries({
      queryKey: [getSearchHistoryOptions().queryKey],
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
