// The production build injects a content-addressed shell and corpus revision.
const BUILD = self.__REMINDER_BUILD__;
const SHELL_CACHE = 'reminder-shell-' + (BUILD?.version || 'dev');
const TEXT_CACHE = 'reminder-text-' + (BUILD?.textVersion || 'dev');
const RECENT_CACHE = 'reminder-recent-v2';
const INDEXES = ['/api/quran/chapters', '/api/hadith/books', '/api/names'];
const SHELL = BUILD?.assets || [];
const FALLBACK = '/__spa-fallback.html';
let downloading;

async function notify(status) {
  const windows = await self.clients.matchAll({ includeUncontrolled: true });
  windows.forEach(client => client.postMessage({ type: 'OFFLINE_STATUS', ...status }));
}

async function network(request) {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), 10000);
  try { return await fetch(request, { signal: controller.signal }); }
  finally { clearTimeout(timer); }
}

function unavailable() {
  return new Response(JSON.stringify({ message: 'This text has not been saved on this device yet. Connect to the internet to download it.' }), {
    status: 503, headers: { 'Content-Type': 'application/json' },
  });
}

async function savedText(path) {
  const current = await (await caches.open(TEXT_CACHE)).match(path);
  if (current) return current;
  // Keep the last downloaded corpus usable during an interrupted update.
  for (const name of (await caches.keys()).filter(name => name.startsWith('reminder-text-') && name !== TEXT_CACHE)) {
    const saved = await (await caches.open(name)).match(path);
    if (saved) return saved;
  }
}

async function textResponse(path) {
  const cached = await savedText(path);
  if (cached) return cached;
  // Name details are already included in the downloaded names list.
  const name = path.match(/^\/api\/names\/(\d+)$/);
  if (name) {
    const list = await savedText('/api/names');
    const item = list && (await list.json()).find(item => item.number === Number(name[1]));
    if (item) return new Response(JSON.stringify(item), { headers: { 'Content-Type': 'application/json' } });
  }
  try {
    const response = await network(path);
    if (!response.ok) return unavailable();
    await response.clone().json(); // Never store HTML fallback/error responses as text.
    await (await caches.open(TEXT_CACHE)).put(path, response.clone()).catch(() => {});
    return response;
  } catch { return unavailable(); }
}

async function saveAllText() {
  const cache = await caches.open(TEXT_CACHE);
  let completed = 0;
  let total = INDEXES.length;
  async function save(path) {
    let response = await cache.match(path);
    if (!response) {
      response = await network(path);
      if (!response.ok) throw new Error('download');
      await response.clone().json();
      await cache.put(path, response.clone());
    }
    completed++;
    return response;
  }
  try {
    const chapters = await (await save(INDEXES[0])).json();
    const books = await (await save(INDEXES[1])).json();
    await save(INDEXES[2]);
    if (!Array.isArray(chapters) || chapters.length !== 114 || !Array.isArray(books) || !books.length) throw new Error('index');
    const paths = [...chapters.map(ch => '/api/quran/' + ch.number), ...books.map(book => '/api/hadith/' + book.number)];
    total += paths.length;
    await notify({ state: 'saving', completed, total });
    // Bound concurrency and resume from successful entries after interruption.
    let next = 0;
    let failed = false;
    await Promise.all(Array.from({ length: 3 }, async () => {
      while (next < paths.length && !failed) {
        const path = paths[next++];
        try { await save(path); await notify({ state: 'saving', completed, total }); }
        catch { failed = true; }
      }
    }));
    if (failed) throw new Error('download');
    await notify({ state: 'ready', completed, total });
    for (const name of await caches.keys()) {
      if ((name.startsWith('reminder-text-') && name !== TEXT_CACHE) || ['reminder-v1', 'reminder-api-v1'].includes(name)) await caches.delete(name);
    }
  } catch {
    await notify({ state: 'paused', completed, total });
  }
}

function downloadText() {
  if (!downloading) downloading = saveAllText().finally(() => { downloading = undefined; });
  return downloading;
}

self.addEventListener('install', event => {
  if (!BUILD) return; // Production precache is generated after the Vite build.
  event.waitUntil((async () => {
    await (await caches.open(SHELL_CACHE)).addAll(SHELL);
    await self.skipWaiting();
  })());
});

self.addEventListener('activate', event => {
  event.waitUntil((async () => {
    // Keep the preceding shell for tabs still running its JavaScript after an update.
    const shells = (await caches.keys()).filter(name => name.startsWith('reminder-shell-') && name !== SHELL_CACHE);
    for (const name of shells.slice(0, -1)) await caches.delete(name);
    await self.clients.claim();
  })());
});

self.addEventListener('message', event => {
  if (event.data?.type === 'DOWNLOAD_TEXT' && BUILD) event.waitUntil(downloadText());
});

self.addEventListener('fetch', event => {
  const request = event.request;
  const url = new URL(request.url);
  if (!BUILD || request.method !== 'GET' || url.origin !== self.location.origin) return;
  if (/^\/api\/(quran\/(chapters|\d+)|hadith\/(books|\d+)|names(?:\/\d+)?)$/.test(url.pathname)) {
    event.respondWith(textResponse(url.pathname));
    return;
  }
  if (url.pathname.startsWith('/api/')) {
    // Dynamic content remains network-first. Only public daily responses are cached.
    event.respondWith((async () => {
      const cacheable = /^\/api\/(latest|daily(?:\/(index|\d{4}-\d{2}-\d{2}))?)$/.test(url.pathname);
      const cache = await caches.open(RECENT_CACHE);
      try {
        const response = await network(request);
        if (response.ok) {
          if (cacheable) await cache.put(request, response.clone()).catch(() => {});
          return response;
        }
        return (cacheable && await cache.match(request)) || response;
      } catch { return (cacheable && await cache.match(request)) || unavailable(); }
    })());
    return;
  }
  event.respondWith((async () => {
    const cache = await caches.open(SHELL_CACHE);
    if (request.mode === 'navigate') {
      // One version-consistent SPA shell works at every deep link, even unopened routes.
      return await cache.match(FALLBACK) || await network(request);
    }
    const saved = await cache.match(request);
    if (saved) return saved;
    for (const name of (await caches.keys()).filter(name => name.startsWith('reminder-shell-') && name !== SHELL_CACHE)) {
      const previous = await (await caches.open(name)).match(request);
      if (previous) return previous;
    }
    return await network(request);
  })());
});

// Push notification handling
self.addEventListener('push', function(event) {
  var data = {};
  if (event.data) {
    data = event.data.json();
  }
  var pushUrl = (data.data && data.data.url) || '/';
  var title = data.title || 'Reminder';
  var options = {
    body: data.body || '',
    icon: '/icon-192.png',
    badge: '/icon-96.png',
    data: {
      url: pushUrl,
    },
  };
  event.waitUntil(self.registration.showNotification(title, options));
});

self.addEventListener('notificationclick', function(event) {
  event.notification.close();
  event.waitUntil(
    clients.openWindow(event.notification.data.url)
  );
});
