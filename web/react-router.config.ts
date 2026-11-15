import type { Config } from "@react-router/dev/config";

export default {
  // The Go server serves build-time HTML; no Node server is deployed.
  ssr: false,
  // Offline navigation cannot fetch lazy route-manifest patches.
  routeDiscovery: { mode: 'initial' },
  prerender: {
    unstable_concurrency: 8,
    async paths() {
      const { readFile } = await import('node:fs/promises');
      return JSON.parse(await readFile('.generated/paths.json', 'utf8'));
    },
  },
} satisfies Config;
