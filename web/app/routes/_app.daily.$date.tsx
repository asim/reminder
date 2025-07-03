import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';

interface DailyData {
  name?: string;
  hadith?: string;
  verse?: string;
  links?: Record<string, string>;
  updated?: string;
  message?: string;
  date?: string;
  hijri?: string;
  error?: string;
}

export default function DailyByDate() {
  const { date } = useParams<{ date: string }>();
  const [data, setData] = useState<DailyData>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!date) return;
    setLoading(true);
    fetch(`/api/daily`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ date }),
    })
      .then(async (res) => {
        if (!res.ok) {
          setData({ error: 'Not found' });
        } else {
          setData(await res.json());
        }
        setLoading(false);
      })
      .catch(() => {
        setData({ error: 'Error loading data' });
        setLoading(false);
      });
  }, [date]);

  if (!date) {
    return <div className="p-8 text-red-600">No date provided</div>;
  }
  if (loading) {
    return <div className="p-8">Loading...</div>;
  }
  if (data.error) {
    return <div className="p-8 text-red-600">{data.error}</div>;
  }
  return (
    <div className="max-w-2xl mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Daily Reminder for {data.date}</h1>
      <div className="mb-4">
        <div className="font-semibold">Hijri Date:</div>
        <div>{data.hijri}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Verse:</div>
        <div>{data.verse}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Hadith:</div>
        <div>{data.hadith}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Name:</div>
        <div>{data.name}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Message:</div>
        <div>{data.message}</div>
      </div>
      <div className="mb-4">
        <div className="font-semibold">Links:</div>
        <ul>
          {data.links && Object.entries(data.links).map(([key, value]) => (
            <li key={key}>
              <a href={String(value)} className="text-blue-600 underline">{key}</a>
            </li>
          ))}
        </ul>
      </div>
      <div className="text-xs text-gray-500">Updated: {data.updated}</div>
    </div>
  );
}
