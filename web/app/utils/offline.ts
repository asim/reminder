import { queryClient } from '~/utils/query-client';

// Start once per document. A terminated/interrupted worker resumes cached progress
// on the next visit or reconnection; no completion flag can outlive evicted data.
export function startOfflineDownload() {
  if (!('serviceWorker' in navigator) || !import.meta.env.PROD) return () => {};
  let registration: ServiceWorkerRegistration | undefined;
  let lastCheck = 0;
  const checkForUpdate = () => {
    if (!registration || !navigator.onLine || document.visibilityState !== 'visible' || Date.now() - lastCheck < 60000) return;
    lastCheck = Date.now();
    void registration.update().catch(console.error);
  };
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
  window.addEventListener('online', checkForUpdate);
  document.addEventListener('visibilitychange', checkForUpdate);
  navigator.serviceWorker.register('/reminder.js', { updateViaCache: 'none' })
    .then(reg => { registration = reg; checkForUpdate(); return navigator.serviceWorker.ready; }).then(download).catch(console.error);
  return () => {
    navigator.serviceWorker.removeEventListener('message', updated);
    navigator.serviceWorker.removeEventListener('controllerchange', download);
    window.removeEventListener('online', download);
    window.removeEventListener('online', checkForUpdate);
    document.removeEventListener('visibilitychange', checkForUpdate);
  };
}
