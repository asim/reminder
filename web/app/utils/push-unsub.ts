// Utility for registering service worker, subscribing and unsubscribing to push
export async function registerServiceWorker() {
  if ('serviceWorker' in navigator) {
    return navigator.serviceWorker.register('/service-worker.js');
  }
  throw new Error('Service workers are not supported');
}

export async function subscribeUserToPush(publicKey: string) {
  // ...existing code...
}

export async function unsubscribeUserFromPush() {
  if (!('serviceWorker' in navigator)) throw new Error('Service workers not supported');
  const reg = await navigator.serviceWorker.ready;
  const sub = await reg.pushManager.getSubscription();
  if (sub) {
    await fetch('/api/push/unsubscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ endpoint: sub.endpoint }),
    });
    await sub.unsubscribe();
    return true;
  }
  return false;
}

export function urlBase64ToUint8Array(base64String: string) {
  // ...existing code...
}
