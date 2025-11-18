import { useQuery } from '@tanstack/react-query';
import { Loader2 } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { getSearchHistoryOptions, searchOptions, type SearchResponseType } from '~/queries/search';
import {
  cacheSearchResult,
  getAllCachedSearches,
  type CachedSearch
} from '~/utils/search-cache';

export function meta() {
  return [
    { title: 'Search - Reminder' },
    {
      property: 'og:title',
      content: 'Search - Reminder',
    },
    {
      name: 'description',
      content:
        'Search the Quran, Hadith, and more through our app and API',
    },
  ];
}

export default function SearchIndex() {
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [expandedRefs, setExpandedRefs] = useState<Record<number, boolean>>({});
  const [showReferences, setShowReferences] = useState(false);
  const [cachedSearches, setCachedSearches] = useState<CachedSearch[]>([]);
  const [expandedHistoryRefs, setExpandedHistoryRefs] = useState<Record<string, boolean>>({});
  const [expandedHistoryRefItems, setExpandedHistoryRefItems] = useState<Record<string, boolean>>({});

  // Load cached searches on mount
  useEffect(() => {
    setCachedSearches(getAllCachedSearches());
  }, []);

  const { data: historyData, refetch: refetchHistory } = useQuery(
    getSearchHistoryOptions()
  );

  const { data: searchResults, isLoading } = useQuery({
    ...searchOptions(submittedQuery),
    enabled: !!submittedQuery,
  });

  // Cache search results when received
  useEffect(() => {
    if (searchResults && submittedQuery) {
      cacheSearchResult(submittedQuery, searchResults);
      setCachedSearches(getAllCachedSearches());
    }
  }, [searchResults, submittedQuery]);

  const handleSubmit = (e: React.FormEvent) => {
    console.log('handleSubmit called');
    e.preventDefault();
    if (!query.trim()) {
      return;
    }

    setSubmittedQuery(query);
    setShowReferences(false);
    setQuery('');
    refetchHistory();
  };

  // Create merged history with cached data
  const mergedHistory = useMemo(() => {
    if (!historyData?.history || historyData.history.length === 0) {
      // No server history, use only cached searches
      return cachedSearches.map(cached => ({
        type: 'cached' as const,
        result: cached.result,
      }));
    }

    // Build a map of cached searches by query for quick lookup
    const cachedMap = new Map<string, CachedSearch>();
    cachedSearches.forEach(cached => {
      cachedMap.set(cached.query.toLowerCase(), cached);
    });

    // Merge server history with cached data
    const merged: Array<{
      type: 'cached' | 'server';
      query?: string;
      answer?: string;
      result?: SearchResponseType;
    }> = [];

    for (let i = 0; i < historyData.history.length; i += 2) {
      const question = historyData.history[i];
      const answer = historyData.history[i + 1];

      if (!question || !answer) continue;

      // Try to find cached version with references
      const cached = cachedMap.get(question.toLowerCase());

      if (cached) {
        merged.push({
          type: 'cached',
          result: cached.result,
        });
      } else {
        merged.push({
          type: 'server',
          query: question,
          answer: answer,
        });
      }
    }

    return merged;
  }, [historyData, cachedSearches]);

  const toggleReference = (index: number) => {
    setExpandedRefs((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const toggleReferencesSection = () => {
    setShowReferences(!showReferences);
  };

  const toggleHistoryReferences = (historyKey: string) => {
    setExpandedHistoryRefs((prev) => ({
      ...prev,
      [historyKey]: !prev[historyKey],
    }));
  };

  const toggleHistoryReference = (historyKey: string, refIndex: number) => {
    const key = `${historyKey}-${refIndex}`;
    setExpandedHistoryRefItems((prev) => ({
      ...prev,
      [key]: !prev[key],
    }));
  };

  return (
    <div className='w-full flex flex-col flex-1 overflow-y-auto'>
      <div className='flex flex-col flex-1 w-full p-4 sm:p-6 lg:p-8 max-w-3xl mx-auto'>
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
          Search
        </h1>
        <div className="text-sm sm:text-base text-gray-700 mb-2">
          Seek knowledge from the Quran, Hadith and names of Allah
        </div>
        <form onSubmit={handleSubmit} method="post" className='mb-4 sm:mb-6'>
          <div className='relative flex'>
            <input
              className='w-full p-2 sm:p-3 border border-gray-300 rounded-md text-sm sm:text-base disabled:opacity-90'
              placeholder={isLoading ? 'Seeking...' : 'Ask a question'}
              value={query}
              onChange={(e) => {
                console.log('Input changed:', e.target.value);
                setQuery(e.target.value);
              }}
              autoComplete='off'
              disabled={isLoading}
              autoFocus
            />
            <button
              type='submit'
              className='ml-2 px-3 py-2 rounded bg-black text-white text-xs sm:text-sm hover:bg-gray-800 disabled:opacity-60'
              disabled={isLoading || !query.trim()}
            >
              Search
            </button>
            {isLoading && (
              <div className='absolute right-3 top-1/2 -translate-y-1/2'>
                <Loader2 className='animate-spin h-4 w-4 sm:h-5 sm:w-5 text-gray-400' />
              </div>
            )}
          </div>
        </form>

        {isLoading && (
          <div className='mb-3 sm:mb-4 text-sm sm:text-base'>Seeking...</div>
        )}

        {searchResults && (
          <div className='mb-6 sm:mb-8 border border-gray-200 rounded-md p-3 sm:p-4'>
            <div className='mb-2 sm:mb-3 text-lg sm:text-xl font-medium'>
              {searchResults.q}
            </div>
            <div
              className='mb-3 sm:mb-4 text-sm sm:text-base prose prose-sm sm:prose-base max-w-none'
              dangerouslySetInnerHTML={{ __html: searchResults.answer }}
            />

            <div className='mt-3 sm:mt-4 border-t pt-3'>
              <button
                onClick={toggleReferencesSection}
                className='text-xs sm:text-sm font-medium text-gray-700 hover:text-black flex items-center gap-1'
              >
                <span>{showReferences ? '▼' : '▶'}</span>
                <span>References ({searchResults.references.length})</span>
              </button>

              {showReferences && (
                <div className='mt-3 space-y-3'>
                  {searchResults.references.map((ref, index) => (
                    <div key={index} className='border border-gray-200 rounded-lg overflow-hidden'>
                      <div
                        className='bg-gray-50 px-3 py-2 cursor-pointer hover:bg-gray-100 flex items-start justify-between gap-2'
                        onClick={() => toggleReference(index)}
                      >
                        <div className='flex-1 min-w-0'>
                          <div className='text-xs sm:text-sm font-medium text-gray-900 truncate'>
                            {ref.metadata.type || 'Reference'} {ref.metadata.chapter && `- Chapter ${ref.metadata.chapter}`}
                            {ref.metadata.verse && `:${ref.metadata.verse}`}
                            {ref.metadata.hadith && `- Hadith ${ref.metadata.hadith}`}
                          </div>
                          <div className='text-xs text-gray-500 mt-1'>
                            {ref.text.substring(0, 80)}...
                          </div>
                        </div>
                        <div className='flex items-center gap-2 flex-shrink-0'>
                          <span className='text-xs text-gray-500'>
                            {(ref.score * 100).toFixed(0)}%
                          </span>
                          <span className='text-gray-400'>
                            {expandedRefs[index] ? '▼' : '▶'}
                          </span>
                        </div>
                      </div>

                      {expandedRefs[index] && (
                        <div className='px-3 py-3 bg-white space-y-2'>
                          <div className='text-xs sm:text-sm text-gray-800 leading-relaxed'>
                            {ref.text}
                          </div>

                          {Object.keys(ref.metadata).length > 0 && (
                            <div className='pt-2 border-t border-gray-100'>
                              <div className='text-xs font-medium text-gray-600 mb-1'>Source Information:</div>
                              <div className='grid grid-cols-2 gap-x-3 gap-y-1 text-xs text-gray-600'>
                                {Object.entries(ref.metadata).map(([key, value]) => (
                                  <div key={key} className='flex gap-1'>
                                    <span className='font-medium capitalize'>{key}:</span>
                                    <span>{value}</span>
                                  </div>
                                ))}
                              </div>
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {mergedHistory && mergedHistory.length > 0 && (
          <div>
            <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>
              Recent Searches
            </h2>
            <div className='space-y-6'>
              {mergedHistory.map((item, idx) => {
                // Handle cached search objects
                if (item.type === 'cached' && item.result) {
                  const result = item.result;
                  const historyKey = `history-${idx}`;
                  return (
                    <div key={idx} className='border border-gray-200 rounded-md p-3 sm:p-4'>
                      <div className='mb-2 font-semibold bg-yellow-100 px-2 py-1 sm:py-2'>
                        {result.q}
                      </div>
                      <div
                        className='mb-3 text-sm sm:text-base prose prose-sm sm:prose-base max-w-none'
                        dangerouslySetInnerHTML={{ __html: result.answer }}
                      />
                      {result.references && result.references.length > 0 && (
                        <div className='mt-3 border-t pt-3'>
                          <button
                            onClick={() => toggleHistoryReferences(historyKey)}
                            className='text-xs sm:text-sm font-medium text-gray-700 hover:text-black flex items-center gap-1'
                          >
                            <span>{expandedHistoryRefs[historyKey] ? '▼' : '▶'}</span>
                            <span>References ({result.references.length})</span>
                          </button>
                          
                          {expandedHistoryRefs[historyKey] && (
                            <div className='mt-3 space-y-3'>
                              {result.references.map((ref, refIdx) => {
                                const refKey = `${historyKey}-${refIdx}`;
                                return (
                                  <div key={refIdx} className='border border-gray-200 rounded-lg overflow-hidden'>
                                    <div
                                      className='bg-gray-50 px-3 py-2 cursor-pointer hover:bg-gray-100 flex items-start justify-between gap-2'
                                      onClick={() => toggleHistoryReference(historyKey, refIdx)}
                                    >
                                      <div className='flex-1 min-w-0'>
                                        <div className='text-xs sm:text-sm font-medium text-gray-900 truncate'>
                                          {ref.metadata.type || 'Reference'} {ref.metadata.chapter && `- Chapter ${ref.metadata.chapter}`}
                                          {ref.metadata.verse && `:${ref.metadata.verse}`}
                                          {ref.metadata.hadith && `- Hadith ${ref.metadata.hadith}`}
                                        </div>
                                        <div className='text-xs text-gray-500 mt-1'>
                                          {ref.text.substring(0, 80)}...
                                        </div>
                                      </div>
                                      <div className='flex items-center gap-2 flex-shrink-0'>
                                        <span className='text-xs text-gray-500'>
                                          {(ref.score * 100).toFixed(0)}%
                                        </span>
                                        <span className='text-gray-400'>
                                          {expandedHistoryRefItems[refKey] ? '▼' : '▶'}
                                        </span>
                                      </div>
                                    </div>

                                    {expandedHistoryRefItems[refKey] && (
                                      <div className='px-3 py-3 bg-white space-y-2'>
                                        <div className='text-xs sm:text-sm text-gray-800 leading-relaxed'>
                                          {ref.text}
                                        </div>

                                        {Object.keys(ref.metadata).length > 0 && (
                                          <div className='pt-2 border-t border-gray-100'>
                                            <div className='text-xs font-medium text-gray-600 mb-1'>Source Information:</div>
                                            <div className='grid grid-cols-2 gap-x-3 gap-y-1 text-xs text-gray-600'>
                                              {Object.entries(ref.metadata).map(([key, value]) => (
                                                <div key={key} className='flex gap-1'>
                                                  <span className='font-medium capitalize'>{key}:</span>
                                                  <span>{value}</span>
                                                </div>
                                              ))}
                                            </div>
                                          </div>
                                        )}
                                      </div>
                                    )}
                                  </div>
                                );
                              })}
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  );
                }

                // Handle server history without cached references
                if (item.type === 'server' && item.query && item.answer) {
                  return (
                    <div key={idx} className='border border-gray-200 rounded-md p-3 sm:p-4'>
                      <div className='mb-2 font-semibold bg-yellow-100 px-2 py-1 sm:py-2'>
                        {item.query}
                      </div>
                      <div
                        className='text-sm sm:text-base prose prose-sm sm:prose-base max-w-none'
                        dangerouslySetInnerHTML={{ __html: item.answer }}
                      />
                    </div>
                  );
                }

                return null;
              })}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
