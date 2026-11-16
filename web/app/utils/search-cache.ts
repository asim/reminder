import type { SearchResponseType } from '~/queries/search';

const SEARCH_CACHE_KEY = 'reminder_search_cache';
const NAMES_REVISION = '2026-11-15-muslim-sqlite';

function identity(query: string, source = 'all', mode = 'search') {
    return JSON.stringify([query.trim().toLowerCase(), source || 'all', mode]);
}

export type CachedSearch = {
    query: string;
    result: SearchResponseType;
    timestamp: number;
};

export type SearchCache = {
    searches: CachedSearch[];
    namesRevision?: string;
};

// Get all cached searches
export function getSearchCache(): SearchCache {
    if (typeof window === 'undefined') {
        return { searches: [] };
    }

    try {
        const cached = localStorage.getItem(SEARCH_CACHE_KEY);
        if (!cached) {
            return { searches: [] };
        }
        const cache: SearchCache = JSON.parse(cached);
        if (cache.namesRevision !== NAMES_REVISION) {
            // New collections change the evidence available to every search.
            cache.searches = [];
            cache.namesRevision = NAMES_REVISION;
            localStorage.setItem(SEARCH_CACHE_KEY, JSON.stringify(cache));
        }
        return cache;
    } catch (error) {
        console.error('Error reading search cache:', error);
        return { searches: [] };
    }
}

// Save a search result to cache
export function cacheSearchResult(query: string, result: SearchResponseType) {
    if (typeof window === 'undefined') return;

    try {
        const cache = getSearchCache();

        // Remove any existing cache for this query
        const filteredSearches = cache.searches.filter(
            (s) => identity(s.query, s.result.source, s.result.mode || (s.result.answer ? 'ask' : 'search')) !== identity(query, result.source, result.mode || (result.answer ? 'ask' : 'search'))
        );

        // Add new search to the beginning
        filteredSearches.unshift({
            query,
            result,
            timestamp: Date.now(),
        });

        // Keep only the last 50 searches
        const trimmedSearches = filteredSearches.slice(0, 50);

        localStorage.setItem(
            SEARCH_CACHE_KEY,
            JSON.stringify({ searches: trimmedSearches, namesRevision: NAMES_REVISION })
        );
    } catch (error) {
        console.error('Error caching search result:', error);
    }
}

// Get a cached search result by query
export function getCachedSearchResult(query: string, source = 'all', mode = 'search'): SearchResponseType | null {
    if (typeof window === 'undefined') return null;

    try {
        const cache = getSearchCache();
        const cached = cache.searches.find(
            (s) => identity(s.query, s.result.source, s.result.mode || (s.result.answer ? 'ask' : 'search')) === identity(query, source, mode)
        );
        return cached ? cached.result : null;
    } catch (error) {
        console.error('Error retrieving cached search:', error);
        return null;
    }
}

// Get all cached searches with full data (for history display)
export function getAllCachedSearches(): CachedSearch[] {
    return getSearchCache().searches;
}

// Clear the entire cache
export function clearSearchCache() {
    if (typeof window === 'undefined') return;

    try {
        localStorage.removeItem(SEARCH_CACHE_KEY);
    } catch (error) {
        console.error('Error clearing search cache:', error);
    }
}

// Restore each mode independently: AI evidence is not an ordinary search result.
export function restoreCachedSearches(
    query: string,
    source: string,
    mode: 'search' | 'ask',
    restore: (mode: 'search' | 'ask', result: SearchResponseType) => void,
) {
    const search = getCachedSearchResult(query, source, 'search');
    if (search) restore('search', search);
    if (mode === 'ask') {
        const summary = getCachedSearchResult(query, source, 'ask');
        if (summary) restore('ask', summary);
    }
}
