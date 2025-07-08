import React, { useState } from 'react';
import { subscribeUserToPush } from '../utils/push';
import { unsubscribeUserFromPush } from '../utils/push-unsub';
// VAPID public key endpoint
const VAPID_PUBLIC_KEY_ENDPOINT = '/api/push/key';
function NotificationButton() {
  const [enabled, setEnabled] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Check subscription status on mount
  React.useEffect(() => {
    if ('serviceWorker' in navigator && 'PushManager' in window) {
      navigator.serviceWorker.getRegistration('/reminder.js').then(reg => {
        if (!reg) return setEnabled(false);
        reg.pushManager.getSubscription().then(sub => {
          setEnabled(!!sub);
        });
      });
    }
  }, []);

  async function handleSubscribe() {
    setLoading(true);
    setError(null);
    try {
      const resp = await fetch(VAPID_PUBLIC_KEY_ENDPOINT);
      if (!resp.ok) throw new Error('Failed to get VAPID key');
      const key = await resp.text();
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
    <div>
      {enabled ? (
        <button
          className="bg-gray-500 text-white px-4 py-2 rounded text-sm cursor-pointer shadow hover:bg-gray-600"
          onClick={handleUnsubscribe}
          disabled={loading}
        >
          Unsubscribe
        </button>
      ) : (
        <button
          className="bg-gray-800 text-white px-4 py-2 rounded text-sm cursor-pointer shadow hover:bg-gray-700"
          onClick={handleSubscribe}
          disabled={loading}
        >
          Subscribe
        </button>
      )}
      {error && <div className="text-red-500 mt-2 text-sm">{error}</div>}
    </div>
  );
}

export default function DailyIndex() {
  return (
    <div className="max-w-4xl mx-auto w-full mb-8 sm:mb-12 flex-grow p-0 lg:p-8 space-y-8">
      <div className="flex items-center justify-between mb-4 sm:mb-6">
        <h1 className="text-2xl sm:text-3xl md:text-4xl font-semibold text-left">
          Daily Reminder
        </h1>
        <div className="ml-4"><NotificationButton /></div>
      </div>
      <div>In the name of Allah, the most beneficent, the most merciful</div>
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">What is the Reminder?</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          The reminder (a word often used to describe the Quran) is an app and API for the Quran, 
          hadith and names of Allah. It is a way to share the message of Islam with everyone in need.
          A resource by which we can renew our intention and work towards the best result in the afterlife 
          inshallah.
        </div>
      </section>
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">Why are we here?</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          This life is a test. We were put here by God to know him and to worship him. To do good deeds and uphold 
          the obligatory acts of faith, prayer, charity, fasting and hajj as demonstrated by our prophet Muhammad (pbuh).
          We are in constant need of this Reminder. Let us internalise the purpose of our creation.
          Read a verse, hadith and name of Allah to reflect, reset and strengthen your intention.
          May we reap the rewards of our actions in this life and the next.
        </div>
      </section>
      <section>
        <h2 className="text-lg sm:text-xl font-medium mb-1 sm:mb-2">Who is Allah?</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Allah is the creator of everything. He is the one and only God.

          <p className="mt-2">Surah Al-Ikhlas (The Sincerity) describes it best:</p>

          <p className="border-l-5 p-4 mt-2 mb-2 italic">
            He is Allah, the One.<br />
            The eternal, the Absolute.<br />
            He did not beget, nor was he begotten.<br />
            And there is none like him. 
            <br /><br />

            <a href="/quran/112" className="text-sm">Quran 112</a>
          </p>
        </div>
      </section>
        <section className='mt-4 sm:mt-4 mb-8 sm:mb-8 bg-gray-50 p-3 sm:p-4 rounded-lg border border-gray-200'>
          <h2 className='text-base sm:text-lg font-medium mb-1 sm:mb-2 text-gray-800'>
            Navigating the Reminder
          </h2>
          <p className='text-sm sm:text-base text-gray-600'>
            Select a daily reminder from the menu which encompasses a verse of the Quran, 
            hadith from sahih al-bukhari and a name of Allah. Read, reflect and reset. 
            If you need further reminders, see hourly updates in the latest tab.
            Once you have read a few, continue strengthening your faith by reading more 
            of the Quran, hadith or names of Allah in our app and making it a daily habit.
            Subscribe to daily notifications using the Subscribe button at the top.
          </p>
        </section>
    </div>
  );
}
