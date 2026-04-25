const CACHE_NAME = 'reminder-v1';
const API_CACHE = 'reminder-api-v1';

// App shell files cached on install
const APP_SHELL = [
  '/',
  '/home',
  '/quran',
  '/hadith',
  '/names',
  '/manifest.webmanifest',
  '/icon-192.png',
  '/icon-96.png',
  '/reminder.png',
  '/fonts/arabic.otf',
];

// API routes to cache for offline access
const API_ROUTES = [
  '/api/quran/chapters',
  '/api/hadith/books',
  '/api/names',
  '/api/latest',
];

// Pre-cache the app shell on install
self.addEventListener('install', function(event) {
  event.waitUntil(
    caches.open(CACHE_NAME).then(function(cache) {
      return cache.addAll(APP_SHELL);
    }).then(function() {
      return self.skipWaiting();
    })
  );
});

// Clean up old caches on activate
self.addEventListener('activate', function(event) {
  event.waitUntil(
    caches.keys().then(function(names) {
      return Promise.all(
        names.filter(function(name) {
          return name !== CACHE_NAME && name !== API_CACHE;
        }).map(function(name) {
          return caches.delete(name);
        })
      );
    }).then(function() {
      return self.clients.claim();
    })
  );
});

// Network-first for API, cache-first for static assets
self.addEventListener('fetch', function(event) {
  var url = new URL(event.request.url);

  // Skip non-GET requests
  if (event.request.method !== 'GET') return;

  // Skip search and push endpoints (dynamic, shouldn't be cached)
  if (url.pathname === '/api/search' || url.pathname.startsWith('/api/push')) return;

  // API requests: network-first, fall back to cache
  if (url.pathname.startsWith('/api/')) {
    event.respondWith(
      caches.open(API_CACHE).then(function(cache) {
        return fetch(event.request).then(function(response) {
          if (response.ok) {
            cache.put(event.request, response.clone());
          }
          return response;
        }).catch(function() {
          return cache.match(event.request);
        });
      })
    );
    return;
  }

  // Static assets (JS, CSS, fonts, images): cache-first
  if (url.pathname.match(/\.(js|css|otf|woff2?|png|jpg|svg|ico)$/)) {
    event.respondWith(
      caches.match(event.request).then(function(cached) {
        return cached || fetch(event.request).then(function(response) {
          if (response.ok) {
            var clone = response.clone();
            caches.open(CACHE_NAME).then(function(cache) {
              cache.put(event.request, clone);
            });
          }
          return response;
        });
      })
    );
    return;
  }

  // HTML navigation requests: network-first, fall back to cached shell
  if (event.request.headers.get('Accept') && event.request.headers.get('Accept').includes('text/html')) {
    event.respondWith(
      fetch(event.request).then(function(response) {
        if (response.ok) {
          var clone = response.clone();
          caches.open(CACHE_NAME).then(function(cache) {
            cache.put(event.request, clone);
          });
        }
        return response;
      }).catch(function() {
        return caches.match(event.request).then(function(cached) {
          return cached || caches.match('/');
        });
      })
    );
    return;
  }
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
