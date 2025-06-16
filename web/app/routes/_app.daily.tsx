import { useQuery } from '@tanstack/react-query';
import { httpGet, httpPost } from '~/utils/http';
import React, { useState } from 'react';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
}

export default function DailyPage() {
  // Use local state to allow direct update from refresh
  const [localData, setLocalData] = useState<DailyResponse | null>(null);
  const { data, isLoading, error, isFetching } = useQuery<DailyResponse>({
    queryKey: ['daily'],
    queryFn: async () => httpGet<DailyResponse>('/api/daily'),
  });

  // Refresh handler
  const [refreshing, setRefreshing] = useState(false);
  const handleRefresh = async () => {
    setRefreshing(true);
    try {
      const refreshed = await httpPost<DailyResponse>('/api/daily/refresh', {});
      setLocalData(refreshed);
    } catch (e) {
      // Optionally handle error
    } finally {
      setRefreshing(false);
    }
  };

  const displayData = localData || data;

  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <div className="flex items-center justify-between mb-4 sm:mb-6">
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold text-left'>
          Daily Reminder
        </h1>
        <button
          className="ml-4 px-4 py-2 bg-blue-600 text-white rounded shadow hover:bg-blue-700 transition disabled:opacity-50"
          onClick={handleRefresh}
          disabled={isFetching || refreshing}
        >
          {(isFetching || refreshing) ? 'Refreshing...' : 'Refresh'}
        </button>
      </div>
      {isLoading && <p className="text-center">Loading...</p>}
      {error && <p className="text-center text-red-500">Failed to load daily reminder.</p>}
      {displayData && (
        <div className="space-y-8">
          <section>
            <h2 className="text-lg font-semibold mb-2">Quran Verse</h2>
            <div className="bg-blue-50 rounded p-4 text-base shadow">
              <a href={displayData.links['verse']}>{displayData.verse}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Hadith</h2>
            <div className="bg-green-50 rounded p-4 text-base shadow">
              <a href={displayData.links['hadith']}>{displayData.hadith}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
            <div className="bg-yellow-50 rounded p-4 text-base shadow">
              <a href={displayData.links['name']}>{displayData.name}</a>
            </div>
          </section>
          <section>
            <div>Updated {displayData.updated}</div>
          </section>
        </div>
      )}
    </div>
  );
}
