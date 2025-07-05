import { useSuspenseQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import React from 'react';


interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links?: Record<string, string>;
  updated: string;
  message: string;
}

export default function DailyIndex() {
  const { data } = useSuspenseQuery<DailyResponse>({
    queryKey: ['daily'],
    queryFn: async () => httpGet<DailyResponse>('/api/daily'),
  });
  // Defensive: check for required fields in data
  const { verse, hadith, name, links: rawLinks, updated, message } = data;
  const links = rawLinks || {};

  return (
    <div className="max-w-4xl mx-auto w-full mb-8 sm:mb-12 flex-grow p-0 lg:p-8 space-y-8">
      <h1 className="text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left">
        Daily Reminder
      </h1>
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">{message}</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention
        </div>
        <img src="/reflect.jpg" className="rounded"/>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Verse</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A verse from the Quran</div>
        <div className="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
          {links.verse ? <a href={links.verse}>{verse}</a> : verse}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Hadith</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A hadith from sahih bukhari</div>
        <div className="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
          {links.hadith ? <a href={links.hadith}>{hadith}</a> : hadith}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Name of Allah</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A beautiful name from the 99 names of Allah</div>
        <div className="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
          {links.name ? <a href={links.name}>{name}</a> : name}
        </div>
      </section>
      <section>
        <div className="mt-2">Updated {updated}</div>
      </section>
    </div>
  );
}