import React from 'react';
import { useParams } from 'react-router';
import { useQuery } from '@tanstack/react-query';
import { httpPost } from '~/utils/http';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
  date: string;
  hijri?: string;
}

export default function DailyByDate() {
  const { date } = useParams<{ date: string }>();
  const { data, isLoading, error } = useQuery<DailyResponse>({
    queryKey: ['daily-by-date', date],
    queryFn: async () => {
      if (!date) throw new Error('No date provided');
      return httpPost<DailyResponse>('/api/daily', { date });
    },
    enabled: !!date,
  });


  if (!date) return <div className="p-4 text-red-500">No date provided</div>;
  if (isLoading) return <div className="p-4">Loading...</div>;
  if (error || !data) return <div className="p-4 text-red-500">Failed to load daily reminder for {date}.</div>;

  // Defensive: extract fields, but render whatever is present
  const { verse, hadith, name, links: rawLinks, updated, message, hijri: hijriDate } = data;
  const links = rawLinks || {};

  return (
    <div className="space-y-8">
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">{data.message}</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention
        </div>
        <img src="/reflect.jpg" className="rounded mb-2" />
        <div className="text-xs text-gray-500 mb-2">{data.hijri || data.date}</div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Verse</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A verse from the Quran</div>
        <div className="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
          {links.verse ? <a href={links.verse}>{verse}</a> : verse}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Hadith</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A hadith from sahih bukhari</div>
        <div className="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
          {links.hadith ? <a href={links.hadith}>{hadith}</a> : hadith}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A beautiful name from the 99 names of Allah</div>
        <div className="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
          {links.name ? <a href={links.name}>{name}</a> : name}
        </div>
      </section>
      <section>
        <div>Updated {data.updated}</div>
      </section>
    </div>
  );
}
