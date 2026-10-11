import type { Config } from "@react-router/dev/config";

export default {
  // Config options...
  // Server-side render by default, to enable SPA mode set this to `false`
  ssr: false,
  // Offline navigation cannot fetch lazy route-manifest patches.
  routeDiscovery: { mode: 'initial' },
  prerender: [
    '/',
    '/quran',
    '/hadith',
    '/names',
    '/api',
  ]
} satisfies Config;
