// Start once per document. A terminated/interrupted worker resumes cached progress
// on the next visit or reconnection; no completion flag can outlive evicted data.
export function startOfflineDownload() {
  if (!('serviceWorker' in navigator) || !import.meta.env.PROD) return () => {};
  const download = () => {
    if (navigator.onLine) navigator.serviceWorker.controller?.postMessage({ type: 'DOWNLOAD_TEXT' });
  };
  navigator.serviceWorker.addEventListener('controllerchange', download);
  window.addEventListener('online', download);
  navigator.serviceWorker.register('/reminder.js', { updateViaCache: 'none' })
    .then(() => navigator.serviceWorker.ready).then(download).catch(console.error);
  return () => {
    navigator.serviceWorker.removeEventListener('controllerchange', download);
    window.removeEventListener('online', download);
  };
}
