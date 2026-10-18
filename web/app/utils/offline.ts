import { queryClient } from '~/utils/query-client';

// Start once per document. A terminated/interrupted worker resumes cached progress
// on the next visit or reconnection; no completion flag can outlive evicted data.
export function startOfflineDownload() {
  if (!('serviceWorker' in navigator) || !import.meta.env.PROD) return () => {};
  const download = () => {
    if (navigator.onLine) navigator.serviceWorker.controller?.postMessage({ type: 'DOWNLOAD_TEXT' });
  };
  const updated = (event: MessageEvent) => {
    if (event.data?.type !== 'OFFLINE_STATUS' || event.data.state !== 'ready') return;
    // A page may have loaded the previous corpus while its replacement downloaded.
    void queryClient.invalidateQueries({ predicate: query =>
      ['list-surahs', 'get-chapter', 'list-hadith-books', 'get-hadith-book', 'list-names', 'get-name'].includes(String(query.queryKey[0]))
    });
  };
  navigator.serviceWorker.addEventListener('message', updated);
  navigator.serviceWorker.addEventListener('controllerchange', download);
  window.addEventListener('online', download);
  navigator.serviceWorker.register('/reminder.js', { updateViaCache: 'none' })
    .then(() => navigator.serviceWorker.ready).then(download).catch(console.error);
  return () => {
    navigator.serviceWorker.removeEventListener('message', updated);
    navigator.serviceWorker.removeEventListener('controllerchange', download);
    window.removeEventListener('online', download);
  };
}
