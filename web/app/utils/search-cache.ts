import type { SearchResponseType } from '~/queries/search';

const SEARCH_CACHE_KEY = 'reminder_search_cache';

export type CachedSearch = {
  query: string;
  result: SearchResponseType;
  timestamp: number;
};

export type SearchCache = {
  searches: CachedSearch[];
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
    return JSON.parse(cached);
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
      (s) => s.query.toLowerCase() !== query.toLowerCase()
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
      JSON.stringify({ searches: trimmedSearches })
    );
  } catch (error) {
    console.error('Error caching search result:', error);
  }
}

// Get a cached search result by query
export function getCachedSearchResult(query: string): SearchResponseType | null {
  if (typeof window === 'undefined') return null;

  try {
    const cache = getSearchCache();
    const cached = cache.searches.find(
      (s) => s.query.toLowerCase() === query.toLowerCase()
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
