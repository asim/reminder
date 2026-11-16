import { useModelAvailable } from '~/queries/capabilities';
import { useQuery } from '@tanstack/react-query';
import { Loader2, Search as SearchIcon } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { Link } from 'react-router';
import { queryClient } from '~/utils/query-client';
import {
  getSearchHistoryOptions,
  searchOptions,
  searchSummaryOptions,
  type SearchResponseType,
  type SearchSource,
} from '~/queries/search';
import {
  cacheSearchResult,
  getCachedSearchResult,
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
const sourceOptions: { value: SearchSource; label: string }[] = [
  {value: 'all', label: 'All'}, {value: 'quran', label: 'Quran'},
  {value: 'hadith', label: 'Hadith'}, {value: 'tafsir', label: 'Tafsir'}, {value: 'names', label: 'Names'},
];

type SearchReference = {
  text: string;
  score: number;
  metadata: Record<string, string>;
};

function getResultTitle(meta: Record<string, string>): string {
  const source = meta.source || '';
  if (meta.kind === 'mention_inventory') return 'Quran mention lookup';
  if (source === 'lexicon') return `${meta.title} — ${meta.headword || meta.root}`;
  if (source === 'morphology') return `${meta.title} — ${meta.chapter}:${meta.verse}`;
  if (source === 'quran_translation') return `${meta.translator || meta.edition} — ${meta.chapter}:${meta.verse}`;
  if (source === 'quran') {
    const name = meta.name || `Chapter ${meta.chapter}`;
    return meta.verse ? `${name} ${meta.chapter}:${meta.verse}` : `${name} — Chapter ${meta.chapter}`;
  }
  if ((source === 'bukhari' || source === 'muslim')) {
    const book = meta.book || `Book ${meta.book_num}`;
    return `${meta.collection_name || (source === 'muslim' ? 'Sahih Muslim' : 'Sahih Bukhari')} · ${book} - ${source === 'muslim' ? 'Record' : 'Hadith'} ${meta.number}`;
  }
  if (source === 'names') {
    return meta.meaning || meta.english || 'Name of Allah';
  }
  if (source === 'tafsir') {
    return `${meta.title || "Tafsir"} — ${meta.ayahs || `${meta.chapter}:${meta.verse}`}`;
  }
  return 'Result';
}

function getResultLink(meta: Record<string, string>): string | null {
  const source = meta.source || '';
  if ((source === 'quran' || source === 'quran_translation') && meta.chapter) {
    return meta.verse ? `/quran/${meta.chapter}/${meta.verse}` : `/quran/${meta.chapter}`;
  }
  if ((source === 'bukhari' || source === 'muslim') && meta.book_num) {
    return meta.number
      ? `/hadith${source === 'muslim' ? '/muslim' : ''}/${meta.book_num}#${meta.number}`
      : `/hadith${source === 'muslim' ? '/muslim' : ''}/${meta.book_num}`;
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
  if (source === 'quran_translation') return 'Translation';
  if ((source === 'bukhari' || source === 'muslim')) return 'Hadith';
  if (source === 'names') return 'Names of Allah';
  if (source === 'tafsir') return 'Commentary';
  if (source === 'lexicon') return 'Dictionary';
  if (source === 'morphology') return 'Arabic grammar';
  return 'Source';
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
  const modelAvailable = useModelAvailable();
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [mode, setMode] = useState<SearchMode>('search');
  const [source, setSource] = useState<SearchSource>('all');
  const [submittedSource, setSubmittedSource] = useState<SearchSource>('all');
  const [cachedSearches, setCachedSearches] = useState<CachedSearch[]>([]);

  useEffect(() => {
    setCachedSearches(getAllCachedSearches());
    const params = new URLSearchParams(window.location.search);
    const q = params.get('q') || '';
    const filter = sourceOptions.find(item => item.value === params.get('source'))?.value || 'all';
    const savedMode = params.get('mode') === 'ask' ? 'ask' : 'search';
    if (q) {
      const cached = getCachedSearchResult(q, filter, savedMode);
      if (cached) {
        queryClient.setQueryData(searchOptions(q, filter).queryKey, cached);
        if (savedMode === 'ask') queryClient.setQueryData(searchSummaryOptions(q, filter).queryKey, cached);
      }
      setQuery(q); setSource(filter); setSubmittedSource(filter); setMode(savedMode); setSubmittedQuery(q);
    }
  }, []);

  const { data: historyData, refetch: refetchHistory } = useQuery(
    getSearchHistoryOptions()
  );

  const { data: searchResults, isLoading, error: searchError } = useQuery({
    ...searchOptions(submittedQuery, submittedSource),
    enabled: !!submittedQuery,
  });

  const { data: summaryResults, isLoading: isSummaryLoading, error: summaryError } = useQuery({
    ...searchSummaryOptions(submittedQuery, submittedSource),
    enabled:
      modelAvailable && !!submittedQuery && mode === 'ask' && !!searchResults,
  });

  const mergedResults = useMemo<SearchResponseType | null>(() => {
    if (!searchResults) return null;
    return {
      ...searchResults,
      source: submittedSource,
      mode,
      answer: mode === 'ask' ? summaryResults?.answer ?? '' : '',
      references: mode === 'ask' ? summaryResults?.references ?? searchResults.references : searchResults.references,
    };
  }, [searchResults, summaryResults, mode, submittedSource]);

  useEffect(() => {
    if (mergedResults && submittedQuery && (mode === 'search' || !!summaryResults?.answer)) {
      cacheSearchResult(submittedQuery, mergedResults);
      setCachedSearches(getAllCachedSearches());
      if (mode === 'ask') refetchHistory();
    }
  }, [mergedResults, submittedQuery, mode, summaryResults, refetchHistory]);

  const doSubmit = (chosenMode: SearchMode) => {
    if (!query.trim()) return;
    const nextMode = chosenMode === 'ask' && !modelAvailable ? 'search' : chosenMode;
    setMode(nextMode); setSubmittedSource(source); setSubmittedQuery(query.trim());
    const params = new URLSearchParams({q: query.trim(), source, mode: nextMode});
    window.history.replaceState(null, '', `${window.location.pathname}?${params}`);
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    doSubmit('search');
  };

  const mergedHistory = useMemo(() => {
    const merged: Array<{
      type: 'cached' | 'server'; query?: string; answer?: string; result?: SearchResponseType;
    }> = cachedSearches.map(cached => ({type: 'cached', result: cached.result}));
    const cachedQueries = new Set(cachedSearches.filter(cached => !!cached.result.answer).map(cached => cached.query.toLowerCase()));
    const history = historyData?.history || [];
    for (let i = 0; i < history.length; i += 2) {
      const question = history[i], answer = history[i + 1];
      if (question && answer && !cachedQueries.has(question.toLowerCase())) {
        merged.push({type: 'server', query: question, answer});
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
          Seek knowledge from the Quran, Hadith, Tafsir and Names of Allah
        </div>
        <form onSubmit={handleSubmit} className='mb-4 sm:mb-6'>
          <div className='relative flex'>
            <input
              className='w-full p-2 sm:p-3 border border-gray-300 rounded-md text-sm sm:text-base disabled:opacity-90'
              placeholder={isLoading ? 'Seeking...' : modelAvailable ? 'Search or ask a question...' : 'Search...'}
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
          <div className='flex flex-wrap gap-2 mt-3' role='group' aria-label='Search sources'>
            {sourceOptions.map(item => <button key={item.value} type='button'
              aria-pressed={source === item.value} onClick={() => setSource(item.value)}
              className={`rounded-full border px-3 py-1 text-sm ${source === item.value ? 'border-gray-900 bg-gray-900 text-white' : 'border-gray-300 text-gray-700 hover:bg-gray-100'}`}>
              {item.label}
            </button>)}
          </div>
          <div className='flex gap-2 mt-3'>
            <button
              type='submit'
              disabled={isLoading || !query.trim()}
              className='flex items-center gap-1.5 px-3 py-1.5 rounded-md bg-black text-white text-xs sm:text-sm hover:bg-gray-800 disabled:opacity-60 transition-colors'
            >
              <SearchIcon className='size-3.5' />
              Search
            </button>
            {modelAvailable && <button
              type='button'
              onClick={() => doSubmit('ask')}
              disabled={isLoading || !query.trim()}
              className='flex items-center gap-1.5 px-3 py-1.5 rounded-md border border-gray-300 text-xs sm:text-sm hover:bg-gray-100 disabled:opacity-60 transition-colors'
            >
              Ask AI
            </button>}
          </div>
        </form>

        {isLoading && (
          <div className='mb-3 sm:mb-4 text-sm sm:text-base text-gray-500'>
            Searching...
          </div>
        )}

        {searchError && <p className='mb-4 text-sm' role='alert'>Search needs a connection. You can still read downloaded texts from Quran, Hadith and Names.</p>}

        {/* Search mode: show results as cards */}
        {mergedResults && mode === 'search' && (
          <div className='mb-6 sm:mb-8'>
            <div className='mb-3 text-sm text-gray-500'>
              {mergedResults.references.length} results for "{mergedResults.q}" · {sourceOptions.find(item => item.value === submittedSource)?.label}
            </div>
            <div className='space-y-3'>
              {mergedResults.references.map((ref, index) => (
                <ResultCard key={index} ref={ref} />
              ))}
            </div>
          </div>
        )}

        {/* Optional summary with references */}
        {mergedResults && mode === 'ask' && (
          <div className='mb-6 sm:mb-8'>
            <div className='border border-gray-200 rounded-md p-3 sm:p-4 mb-4'>
              <div className='mb-2 sm:mb-3 text-lg sm:text-xl font-medium'>
                {mergedResults.q}
              </div>
              {summaryError ? <p className='text-sm'>{summaryError.message || "AI summary unavailable. The references remain available below."}</p> : isSummaryLoading ? (
                <div className='flex items-center gap-2 text-sm sm:text-base text-gray-500'>
                  <Loader2 className='animate-spin h-4 w-4' />
                  <span>Preparing an answer...</span>
                </div>
              ) : !mergedResults.answer ? <p className='text-sm'>No AI explanation is available. Read the references below.</p> : (
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
                    <div key={index}><span className="text-xs text-gray-500">[{index + 1}]</span><ResultCard ref={ref} compact /></div>
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
                      <p className='mb-2 text-xs text-gray-500'>{result.answer ? 'Ask AI' : 'Search'} · {sourceOptions.find(item => item.value === (result.source || 'all'))?.label}</p>
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
                            <details className='text-sm'><summary className='cursor-pointer text-gray-600'>Show {result.references.length - 3} more results</summary><div className='mt-2 space-y-2'>{result.references.slice(3).map((ref, index) => <ResultCard key={index} ref={ref} />)}</div></details>
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
