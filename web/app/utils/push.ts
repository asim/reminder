// Utility for registering service worker and subscribing to push
export async function registerServiceWorker() {
  if ('serviceWorker' in navigator) {
    return navigator.serviceWorker.register('/service-worker.js');
  }
  throw new Error('Service workers are not supported');
}

export async function subscribeUserToPush(publicKey: string) {
  let registration;
  try {
    registration = await registerServiceWorker();
  } catch (err) {
    console.error('Service worker registration failed:', err);
    throw new Error('Service worker registration failed: ' + err);
  }
  let subscription;
  try {
    subscription = await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(publicKey),
    });
  } catch (err) {
    console.error('PushManager.subscribe failed:', err);
    throw new Error('PushManager.subscribe failed: ' + err);
  }
  try {
    const resp = await fetch('/api/push/subscribe', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(subscription),
    });
    if (!resp.ok) {
      const text = await resp.text();
      console.error('Backend /api/push/subscribe failed:', resp.status, text);
      throw new Error('Backend /api/push/subscribe failed: ' + text);
    }
  } catch (err) {
    console.error('Failed to send subscription to backend:', err);
    throw err;
  }
  return subscription;
}

export function urlBase64ToUint8Array(base64String: string) {
  const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding)
    .replace(/-/g, '+')
    .replace(/_/g, '/');
  const rawData = window.atob(base64);
  const outputArray = new Uint8Array(rawData.length);
  for (let i = 0; i < rawData.length; ++i) {
    outputArray[i] = rawData.charCodeAt(i);
  }
  return outputArray;
}
