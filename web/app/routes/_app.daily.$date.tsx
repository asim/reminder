import React from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { getDailyByDateOptions } from '../queries/daily';
import type { DailyData } from '../queries/daily';

export default function DailyByDate() {
  const { date } = useParams<{ date: string }>();
  const {
    data = {},
    isLoading,
    error,
  } = useQuery<DailyData>(date ? getDailyByDateOptions(date) : { queryKey: ['get-daily', 'none'], queryFn: async () => ({ error: 'No date provided' }) });
  const daily = data as DailyData;

  if (!date) {
    return <div className="p-8 text-red-600">No date provided</div>;
  }
  if (isLoading) {
    return <div className="p-8">Loading...</div>;
  }
  if (error || daily.error) {
    return <div className="p-8 text-red-600">{daily.error || 'Error loading data'}</div>;
  }
  return (
    <div className="max-w-2xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Daily Reminder for {daily.date}</h1>
      <div className="mb-4">
        <div className="font-semibold">Hijri Date:</div>
        <div>{daily.hijri}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Verse:</div>
        <div>{daily.verse}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Hadith:</div>
        <div>{daily.hadith}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Name:</div>
        <div>{daily.name}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Message:</div>
        <div>{daily.message}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Links:</div>
        <ul>
          {daily.links && Object.entries(daily.links).map(([key, value]) => (
            <li key={key}>
              <a href={String(value)} className="text-blue-600 underline">{key}</a>
            </li>
          ))}
        </ul>
      </div>
      <div className="text-xs text-gray-500">Updated: {daily.updated}</div>
    </div>
  );
}
