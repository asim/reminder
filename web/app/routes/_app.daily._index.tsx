import { useSuspenseQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import React, { useState } from 'react';
import { subscribeUserToPush } from '../utils/push';
import { unsubscribeUserFromPush } from '../utils/push-unsub';
// VAPID public key endpoint
const VAPID_PUBLIC_KEY_ENDPOINT = '/api/push/key';
function NotificationButton() {
  const [enabled, setEnabled] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  async function handleSubscribe() {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetch(VAPID_PUBLIC_KEY_ENDPOINT);
      if (!resp.ok) throw new Error('Failed to get VAPID key');
      const { key } = await resp.json();
      await subscribeUserToPush(key);
      setEnabled(true);
    } catch (err: any) {
      setError(err.message || 'Failed to enable notifications');
    } finally {
      setLoading(false);
    }
  }

  async function handleUnsubscribe() {
    setLoading(true);
    setError(null);
    try {
      await unsubscribeUserFromPush();
      setEnabled(false);
    } catch (err: any) {
      setError(err.message || 'Failed to disable notifications');
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="my-4">
      {enabled ? (
        <button
          className="bg-red-500 text-white px-4 py-2 rounded shadow hover:bg-red-600"
          onClick={handleUnsubscribe}
          disabled={loading}
        >
          Disable Notifications
        </button>
      ) : (
        <button
          className="bg-blue-600 text-white px-4 py-2 rounded shadow hover:bg-blue-700"
          onClick={handleSubscribe}
          disabled={loading}
        >
          Enable Notifications
        </button>
      )}
      {error && <div className="text-red-500 mt-2 text-sm">{error}</div>}
    </div>
  );
}


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
      <div className="flex items-center justify-between mb-4 sm:mb-6">
        <h1 className="text-2xl sm:text-3xl md:text-4xl font-semibold text-left">
          Daily Reminder
        </h1>
        <div className="ml-4"><NotificationButton /></div>
      </div>
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">{message}</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention
        </div>
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
