import { useQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import React from 'react';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
}

export default function DailyPage() {
  const { data, isLoading, error } = useQuery<DailyResponse>({
    queryKey: ['daily'],
    queryFn: async () => httpGet<DailyResponse>('/api/daily'),
  });

  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        Daily Reminder
      </h1>
      {isLoading && <p className="text-center">Loading...</p>}
      {error && <p className="text-center text-red-500">Failed to load daily reminder.</p>}
      {data && (
        <div className="space-y-8">
          <section>
            <h2 className="text-lg font-semibold mb-2">Quran Verse</h2>
            <div className="bg-blue-50 rounded p-4 text-base shadow">
              {data.verse}
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Hadith</h2>
            <div className="bg-green-50 rounded p-4 text-base shadow">
              {data.hadith}
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
            <div className="bg-yellow-50 rounded p-4 text-base shadow">
              {data.name}
            </div>
          </section>
        </div>
      )}
    </div>
  );
}
