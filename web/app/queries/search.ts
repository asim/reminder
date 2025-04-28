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

export const searchOptions = (query: string) => ({
  queryKey: ['search', query],
  queryFn: async () => {
    if (!query.trim()) {
      return null;
    }

    const data = await httpPost<SearchResponseType>('/api/search', {
      q: query,
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
