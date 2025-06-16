import { useQuery, useQueryClient } from '@tanstack/react-query';
import { httpGet, httpPost } from '~/utils/http';
import React from 'react';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
}

export default function DailyPage() {
  const queryClient = useQueryClient();
  const { data, isLoading, error, isFetching } = useQuery<DailyResponse>({
    queryKey: ['daily'],
    queryFn: async () => httpGet<DailyResponse>('/api/daily'),
  });

  // Refresh handler
  const handleRefresh = async () => {
    try {
      await httpPost('/api/daily/refresh', {});
      queryClient.invalidateQueries(['daily']);
    } catch (e) {
      // Optionally handle error
    }
  };

  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <div className="flex items-center justify-between mb-4 sm:mb-6">
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold text-left'>
          Daily Reminder
        </h1>
        <button
          className="ml-4 px-4 py-2 bg-blue-600 text-white rounded shadow hover:bg-blue-700 transition disabled:opacity-50"
          onClick={handleRefresh}
          disabled={isFetching}
        >
          {isFetching ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      {isLoading && <p className="text-center">Loading...</p>}
      {error && <p className="text-center text-red-500">Failed to load daily reminder.</p>}
      {data && (
        <div className="space-y-8">
          <section>
            <h2 className="text-lg font-semibold mb-2">Quran Verse</h2>
            <div className="bg-blue-50 rounded p-4 text-base shadow">
              <a href={data.links['verse']}>{data.verse}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Hadith</h2>
            <div className="bg-green-50 rounded p-4 text-base shadow">
              <a href={data.links['hadith']}>{data.hadith}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
            <div className="bg-yellow-50 rounded p-4 text-base shadow">
              <a href={data.links['name']}>{data.name}</a>
            </div>
          </section>
          <section>
            <div>Updated {data.updated}</div>
          </section>
        </div>
      )}
    </div>
  );
}
