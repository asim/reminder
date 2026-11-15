# Search discovery and reading HTML

`npm run build` in `web` exports the bundled Go corpora, pre-renders the existing
React reader routes, and emits robots.txt and sitemap.xml. It requires Go as well
as Node, but no running API or model keys. The release workflow sets up Go first.

The build covers all Quran chapters and numbered verses, hadith books, Names of
Allah, and the public landing/index/API pages. Server loaders supply the same
records as the API to React Query before rendering. Client loaders retain API
and offline-cache navigation. Individual verse pages carry just their own verse
rather than embedding the full chapter again. Reading preferences and bookmarks
are applied after hydration so the initial markup matches the generated defaults.

HTML and route data are gzip-compressed on disk to limit the embedded binary's
size. The Go handler serves gzip to clients that accept it, or decompresses it, with correct content
types, including for HEAD requests. Deploy the Go binary built after `make web`;
a generic static host would need equivalent support for these compressed files.

Only known application routes may fall back to the SPA shell. They receive
`X-Robots-Tag: noindex`. Missing reader pages and assets return HTTP 404. Canonical
URLs use https://reminder.dev and omit reading-preference query parameters. The
sitemap contains the rendered reading URLs, not personalised application routes
or duplicate embed URLs.

Validation:

```sh
make web
python3 scripts/test-rendered.py
go test ./...
REMINDER_RENDERED_DIR="$PWD/web/build/client" go test ./app
cd web && npm run typecheck
```

The rendered-content check inspects every generated page without JavaScript,
including verse text, name summaries, hadith content, titles and canonical URLs.
After deployment, check the live robots/sitemap endpoints and representative
reading URLs, then submit the sitemap and inspect URLs in Google Search Console.
These changes improve crawlability; they do not guarantee inclusion or ranking.
