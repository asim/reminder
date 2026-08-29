import { useQuery } from '@tanstack/react-query';
import { Loader2, MessageSquare, Search as SearchIcon } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router';
import {
  getSearchHistoryOptions,
  searchOptions,
  searchSummaryOptions,
  type SearchResponseType,
} from '~/queries/search';
import {
  cacheSearchResult,
  getAllCachedSearches,
  type CachedSearch,
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

type SearchMode = 'search' | 'ask';

type SearchReference = {
  text: string;
  score: number;
  metadata: Record<string, string>;
};

function getResultTitle(meta: Record<string, string>): string {
  const source = meta.source || '';
  if (source === 'quran') {
    const name = meta.name || `Chapter ${meta.chapter}`;
    return `${name} ${meta.chapter}:${meta.verse}`;
  }
  if (source === 'bukhari') {
    const book = meta.book || `Book ${meta.book_num}`;
    return `${book} - Hadith ${meta.number}`;
  }
  if (source === 'names') {
    return meta.meaning || meta.english || 'Name of Allah';
  }
  if (source === 'tafsir') {
    return `Tafsir - ${meta.chapter}:${meta.verse}`;
  }
  return 'Result';
}

function getResultLink(meta: Record<string, string>): string | null {
  const source = meta.source || '';
  if (source === 'quran' && meta.chapter && meta.verse) {
    return `/quran/${meta.chapter}/${meta.verse}`;
  }
  if (source === 'bukhari' && meta.book_num) {
    return meta.number
      ? `/hadith/${meta.book_num}#${meta.number}`
      : `/hadith/${meta.book_num}`;
  }
  if (source === 'names' && meta.english) {
    return null;
  }
  if (source === 'tafsir' && meta.chapter && meta.verse) {
    return `/quran/${meta.chapter}/${meta.verse}`;
  }
  return null;
}

function getSourceLabel(meta: Record<string, string>): string {
  const source = meta.source || '';
  if (source === 'quran') return 'Quran';
  if (source === 'bukhari') return 'Hadith';
  if (source === 'names') return 'Names of Allah';
  if (source === 'tafsir') return 'Commentary';
  return 'Source';
}

function RelevanceBadge({ score }: { score: number }) {
  const pct = Math.round(score * 100);
  const color =
    pct >= 70
      ? 'bg-green-100 text-green-800'
      : pct >= 40
        ? 'bg-yellow-100 text-yellow-800'
        : 'bg-gray-100 text-gray-600';
  return (
    <span className={`text-xs px-1.5 py-0.5 rounded-full font-medium ${color}`}>
      {pct}% match
    </span>
  );
}

function ResultCard({ ref: r, compact }: { ref: SearchReference; compact?: boolean }) {
  const title = getResultTitle(r.metadata);
  const link = getResultLink(r.metadata);
  const sourceLabel = getSourceLabel(r.metadata);

  return (
    <div className='border border-gray-200 rounded-lg p-3 sm:p-4 hover:border-gray-300 transition-colors'>
      <div className='flex items-start justify-between gap-2 mb-1.5'>
        <div className='flex items-center gap-2 min-w-0'>
          <span className='text-xs px-1.5 py-0.5 rounded bg-gray-100 text-gray-600 font-medium flex-shrink-0'>
            {sourceLabel}
          </span>
          {link ? (
            <Link
              to={link}
              className='text-sm sm:text-base font-medium text-blue-700 hover:text-blue-900 hover:underline truncate'
            >
              {title}
            </Link>
          ) : (
            <span className='text-sm sm:text-base font-medium text-gray-900 truncate'>
              {title}
            </span>
          )}
        </div>
        <RelevanceBadge score={r.score} />
      </div>
      <p className='text-sm text-gray-700 leading-relaxed'>
        {compact ? (r.text.length > 200 ? r.text.slice(0, 200) + '...' : r.text) : r.text}
      </p>
      {r.metadata.narrator && (
        <p className='text-xs text-gray-500 mt-1'>Narrated by {r.metadata.narrator}</p>
      )}
    </div>
  );
}

export default function SearchIndex() {
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [mode, setMode] = useState<SearchMode>('search');
  const [cachedSearches, setCachedSearches] = useState<CachedSearch[]>([]);

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

  const { data: summaryResults, isLoading: isSummaryLoading } = useQuery({
    ...searchSummaryOptions(submittedQuery),
    enabled: !!submittedQuery && mode === 'ask' && !!searchResults,
  });

  const mergedResults = useMemo<SearchResponseType | null>(() => {
    if (!searchResults) return null;
    return {
      ...searchResults,
      answer: summaryResults?.answer ?? '',
    };
  }, [searchResults, summaryResults]);

  useEffect(() => {
    if (mergedResults && submittedQuery && mergedResults.answer) {
      cacheSearchResult(submittedQuery, mergedResults);
      setCachedSearches(getAllCachedSearches());
      refetchHistory();
    }
  }, [mergedResults, submittedQuery, refetchHistory]);

  const doSubmit = (chosenMode: SearchMode) => {
    if (!query.trim()) return;
    setMode(chosenMode);
    setSubmittedQuery(query);
    setQuery('');
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    doSubmit(mode);
  };

  const mergedHistory = useMemo(() => {
    if (!historyData?.history || historyData.history.length === 0) {
      return cachedSearches.map((cached) => ({
        type: 'cached' as const,
        result: cached.result,
      }));
    }

    const cachedMap = new Map<string, CachedSearch>();
    cachedSearches.forEach((cached) => {
      cachedMap.set(cached.query.toLowerCase(), cached);
    });

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

      const cached = cachedMap.get(question.toLowerCase());
      if (cached) {
        merged.push({ type: 'cached', result: cached.result });
      } else {
        merged.push({ type: 'server', query: question, answer });
      }
    }

    return merged;
  }, [historyData, cachedSearches]);

  return (
    <div className='w-full h-full flex flex-col overflow-y-auto'>
      <div className='flex flex-col w-full p-4 sm:p-6 lg:p-8 max-w-3xl mx-auto'>
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
          Search
        </h1>
        <div className='text-sm sm:text-base text-gray-700 mb-2'>
          Seek knowledge from the Quran, Hadith and Names of Allah
        </div>
        <form onSubmit={handleSubmit} className='mb-4 sm:mb-6'>
          <div className='relative flex'>
            <input
              className='w-full p-2 sm:p-3 border border-gray-300 rounded-md text-sm sm:text-base disabled:opacity-90'
              placeholder={isLoading ? 'Seeking...' : 'Search or ask a question...'}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              autoComplete='off'
              disabled={isLoading}
              autoFocus
            />
            {isLoading && (
              <div className='absolute right-2 top-1/2 -translate-y-1/2'>
                <Loader2 className='animate-spin h-4 w-4 text-gray-400' />
              </div>
            )}
          </div>
          <div className='flex gap-2 mt-2'>
            <button
              type='button'
              onClick={() => doSubmit('search')}
              disabled={isLoading || !query.trim()}
              className='flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-black text-white text-xs sm:text-sm hover:bg-gray-800 disabled:opacity-60 transition-colors'
            >
              <SearchIcon className='size-3.5' />
              Search
            </button>
            <button
              type='button'
              onClick={() => doSubmit('ask')}
              disabled={isLoading || !query.trim()}
              className='flex items-center gap-1.5 px-3 py-1.5 rounded-md border border-gray-300 text-gray-700 text-xs sm:text-sm hover:bg-gray-100 disabled:opacity-60 transition-colors'
            >
              <MessageSquare className='size-3.5' />
              Ask AI
            </button>
          </div>
        </form>

        {isLoading && (
          <div className='mb-3 sm:mb-4 text-sm sm:text-base text-gray-500'>
            Searching...
          </div>
        )}

        {/* Search mode: show results as cards */}
        {mergedResults && mode === 'search' && (
          <div className='mb-6 sm:mb-8'>
            <div className='mb-3 text-sm text-gray-500'>
              {mergedResults.references.length} results for "{mergedResults.q}"
            </div>
            <div className='space-y-3'>
              {mergedResults.references.map((ref, index) => (
                <ResultCard key={index} ref={ref} />
              ))}
            </div>
          </div>
        )}

        {/* Ask AI mode: show summary + references */}
        {mergedResults && mode === 'ask' && (
          <div className='mb-6 sm:mb-8'>
            <div className='border border-gray-200 rounded-md p-3 sm:p-4 mb-4'>
              <div className='mb-2 sm:mb-3 text-lg sm:text-xl font-medium'>
                {mergedResults.q}
              </div>
              {isSummaryLoading || !mergedResults.answer ? (
                <div className='flex items-center gap-2 text-sm sm:text-base text-gray-500'>
                  <Loader2 className='animate-spin h-4 w-4' />
                  <span>Generating AI summary...</span>
                </div>
              ) : (
                <div
                  className='text-sm sm:text-base prose prose-sm sm:prose-base max-w-none'
                  dangerouslySetInnerHTML={{ __html: mergedResults.answer }}
                />
              )}
            </div>

            {mergedResults.references.length > 0 && (
              <div>
                <div className='text-sm font-medium text-gray-600 mb-2'>
                  Sources ({mergedResults.references.length})
                </div>
                <div className='space-y-2'>
                  {mergedResults.references.map((ref, index) => (
                    <ResultCard key={index} ref={ref} compact />
                  ))}
                </div>
              </div>
            )}
          </div>
        )}

        {mergedHistory && mergedHistory.length > 0 && (
          <div>
            <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>
              Recent
            </h2>
            <div className='space-y-4'>
              {mergedHistory.map((item, idx) => {
                if (item.type === 'cached' && item.result) {
                  const result = item.result;
                  return (
                    <div
                      key={idx}
                      className='border border-gray-200 rounded-md p-3 sm:p-4'
                    >
                      <div className='mb-2 font-semibold text-sm sm:text-base'>
                        {result.q}
                      </div>
                      {result.answer && (
                        <div
                          className='mb-3 text-sm prose prose-sm max-w-none'
                          dangerouslySetInnerHTML={{ __html: result.answer }}
                        />
                      )}
                      {result.references && result.references.length > 0 && (
                        <div className='space-y-2'>
                          {result.references.slice(0, 3).map((ref, refIdx) => (
                            <ResultCard key={refIdx} ref={ref} compact />
                          ))}
                          {result.references.length > 3 && (
                            <div className='text-xs text-gray-500'>
                              +{result.references.length - 3} more results
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  );
                }

                if (item.type === 'server' && item.query && item.answer) {
                  return (
                    <div
                      key={idx}
                      className='border border-gray-200 rounded-md p-3 sm:p-4'
                    >
                      <div className='mb-2 font-semibold text-sm sm:text-base'>
                        {item.query}
                      </div>
                      <div
                        className='text-sm prose prose-sm max-w-none'
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
