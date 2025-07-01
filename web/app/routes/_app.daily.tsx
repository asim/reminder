import { useQuery } from '@tanstack/react-query';
import { httpGet, httpPost } from '~/utils/http';
import React, { useState } from 'react';
import { subscribeUserToPush } from '../utils/push';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
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

  const [pushStatus, setPushStatus] = useState<'idle' | 'success' | 'error'>('idle');
  const [pushMsg, setPushMsg] = useState('');

  async function handlePushSubscribe() {
    setPushStatus('idle');
    setPushMsg('');
    try {
      // Fetch VAPID public key from backend
      const resp = await fetch('/api/push/vapidPublicKey');
      if (!resp.ok) throw new Error('Failed to fetch VAPID key');
      const publicKey = await resp.text();
      await subscribeUserToPush(publicKey);
      setPushStatus('success');
      setPushMsg('Subscribed to daily notifications!');
    } catch (e) {
      setPushStatus('error');
      setPushMsg('Failed to subscribe to notifications.');
    }
  }

  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <div className="flex items-center justify-between mb-4 sm:mb-6">
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold text-left'>
          Daily Reminder
        </h1>
        <div className='flex gap-2'>
          <button
            className="px-2 py-1 text-sm bg-black text-white rounded shadow hover:bg-gray-800 transition disabled:opacity-50"
            onClick={handleRefresh}
            disabled={isFetching || refreshing}
          >
            {(isFetching || refreshing) ? 'Refreshing...' : 'Refresh'}
          </button>
        </div>
      </div>
      {isLoading && <p className="text-center">Loading...</p>}
      {error && <p className="text-center text-red-500">Failed to load daily reminder.</p>}
      {displayData && (
        <div className="space-y-8">
          {/* Salam and Hijri date message at the top */}
          <section>
            <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">{displayData.message}</h2>
            <div className="text-sm sm:text-base text-gray-700">
              Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Verse</h2>
            <div className="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
              <a href={displayData.links['verse']}>{displayData.verse}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Hadith</h2>
            <div className="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
              <a href={displayData.links['hadith']}>{displayData.hadith}</a>
            </div>
          </section>
          <section>
            <h2 className="text-lg font-semibold mb-2">Name of Allah</h2>
            <div className="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
              <a href={displayData.links['name']}>{displayData.name}</a>
            </div>
          </section>
          <section>
            <div>Updated {displayData.updated}</div>
          </section>
        </div>
      )}
      <div className='flex gap-2 mt-2'>
        <button
          className='px-2 py-1 text-sm bg-blue-600 text-white rounded shadow hover:bg-blue-700 transition'
          onClick={handlePushSubscribe}
        >
          Enable Daily Push Notifications
        </button>
        {pushMsg && (
          <span className={pushStatus === 'success' ? 'text-green-600' : 'text-red-600'}>{pushMsg}</span>
        )}
      </div>
    </div>
  );
}
