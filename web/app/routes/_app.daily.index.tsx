
import { useQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import React from 'react';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
}

export default function DailyIndex() {
  const { data, isLoading, error } = useQuery<DailyResponse>({
    queryKey: ['daily'],
    queryFn: async () => httpGet<DailyResponse>('/api/daily'),
  });

  if (isLoading) return <div className="p-4">Loading...</div>;
  if (error || !data) return <div className="p-4 text-red-500">Failed to load daily reminder.</div>;

  return (
    <div className="space-y-8">
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">{data.message}</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention
        </div>
        <img src="/reflect.jpg" className="rounded"/>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Verse</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A verse from the Quran</div>
        <div className="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
          <a href={data.links['verse']}>{data.verse}</a>
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Hadith</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A hadith from sahih bukhari</div>
        <div className="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
          <a href={data.links['hadith']}>{data.hadith}</a>
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A beautiful name from the 99 names of Allah</div>
        <div className="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
          <a href={data.links['name']}>{data.name}</a>
        </div>
      </section>
      <section>
        <div>Updated {data.updated}</div>
      </section>
    </div>
  );
}
