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
    <div className="max-w-2xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6 text-center">Daily Reminder</h1>
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
