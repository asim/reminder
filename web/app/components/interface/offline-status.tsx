import { useEffect, useState } from 'react';

type Status = { state: 'saving' | 'ready' | 'paused'; completed: number; total: number };

export function OfflineStatus() {
  const [status, setStatus] = useState<Status>();
  const [online, setOnline] = useState(true);
  useEffect(() => {
    const connection = () => setOnline(navigator.onLine);
    const update = (event: MessageEvent) => {
      if (event.data?.type === 'OFFLINE_STATUS') setStatus(event.data);
    };
    connection();
    window.addEventListener('online', connection);
    window.addEventListener('offline', connection);
    navigator.serviceWorker?.addEventListener('message', update);
    return () => {
      window.removeEventListener('online', connection);
      window.removeEventListener('offline', connection);
      navigator.serviceWorker?.removeEventListener('message', update);
    };
  }, []);
  if (online && !status) return null;
  return <div className='px-3 py-1 text-xs text-gray-600 bg-gray-50 border-b border-gray-100' role='status'>
    {!online ? 'Offline · Saved texts are available. Audio and AI need a connection.' :
      status?.state === 'ready' ? 'Quran, hadith and names saved for offline reading.' :
      status?.state === 'saving' ? `Saving texts for offline reading… ${status.completed}/${status.total}` :
      <>Offline download paused. Reconnect or free some storage, then <button className='underline' onClick={() => navigator.serviceWorker.controller?.postMessage({ type: 'DOWNLOAD_TEXT' })}>retry</button>.</>}
  </div>;
}
