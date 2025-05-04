import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { searchOptions, getSearchHistoryOptions } from '~/queries/search';
import { cn } from '~/utils/classname';
import { Loader2 } from 'lucide-react';

export default function SearchIndex() {
  const [query, setQuery] = useState('');
  const [submittedQuery, setSubmittedQuery] = useState('');
  const [expandedRefs, setExpandedRefs] = useState<Record<number, boolean>>({});
  const [showReferences, setShowReferences] = useState(false);

  const { data: historyData, refetch: refetchHistory } = useQuery(
    getSearchHistoryOptions()
  );

  const { data: searchResults, isLoading } = useQuery({
    ...searchOptions(submittedQuery),
    enabled: !!submittedQuery,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!query.trim()) {
      return;
    }

    setSubmittedQuery(query);
    setShowReferences(false);
    setQuery('');
    refetchHistory();
  };

  const toggleReference = (index: number) => {
    setExpandedRefs((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const toggleReferencesSection = () => {
    setShowReferences(!showReferences);
  };

  return (
    <div className='w-full flex flex-col flex-1 overflow-y-auto'>
      <div className='flex flex-col flex-1 w-full p-4 sm:p-6 lg:p-8 max-w-3xl mx-auto'>
        <h1 className='text-xl lg:block hidden sm:text-2xl font-medium mb-2'>Search</h1>

        <form onSubmit={handleSubmit} className='mb-4 sm:mb-6'>
          <div className='relative'>
            <input
              className='w-full p-2 sm:p-3 border border-gray-300 rounded-md text-sm sm:text-base disabled:opacity-90'
              placeholder={isLoading ? 'Seeking...' : 'Ask a question'}
              value={query}
              onChange={(e) => setQuery(e.target.value)}
              autoComplete='off'
              disabled={isLoading}
              autoFocus
            />
            {isLoading && (
              <div className='absolute right-3 top-1/2 -translate-y-1/2'>
                <Loader2 className='animate-spin h-4 w-4 sm:h-5 sm:w-5 text-gray-400' />
              </div>
            )}
          </div>
        </form>

        {isLoading && <div className='mb-3 sm:mb-4 text-sm sm:text-base'>Seeking...</div>}

        {searchResults && (
          <div className='mb-6 sm:mb-8 border border-gray-200 rounded-md p-3 sm:p-4'>
            <div className='mb-2 sm:mb-3 text-lg sm:text-xl font-medium'>{searchResults.q}</div>
            <div
              className='mb-3 sm:mb-4 text-sm sm:text-base'
              dangerouslySetInnerHTML={{ __html: searchResults.answer }}
            />

            <div className='mt-3 sm:mt-4'>
              <button
                onClick={toggleReferencesSection}
                className='text-xs sm:text-sm underline cursor-pointer mb-1'
              >
                {showReferences ? 'Hide References' : 'Show References'}
              </button>

              {showReferences && (
                <div className='mt-1 sm:mt-2'>
                  {searchResults.references.map((ref, index) => (
                    <div key={index} className='mb-3 sm:mb-4'>
                      <div
                        className='text-xs sm:text-sm underline cursor-pointer'
                        onClick={() => toggleReference(index)}
                      >
                        {ref.text.substring(0, 50)}... (Score:{' '}
                        {ref.score.toFixed(2)})
                      </div>

                      {expandedRefs[index] && (
                        <div className='mt-1 ml-2 sm:ml-4 text-xs sm:text-sm'>
                          <div>Text: {ref.text}</div>
                          <div>Metadata: {JSON.stringify(ref.metadata)}</div>
                          <div>Score: {ref.score}</div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}

        {historyData &&
          historyData.history &&
          historyData.history.length > 0 && (
            <div>
              <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>Recent Searches</h2>
              <div className='text-sm sm:text-base'>
                {historyData.history.map((item, idx) => {
                  const isQuestion = !item.startsWith('<p>');
                  const isAnswer = item.startsWith('<p>');

                  return (
                    <div
                      key={idx}
                      className={cn({
                        'font-semibold bg-yellow-100 px-2 py-1 sm:py-2 mb-1 sm:mb-2':
                          isQuestion,
                        'mb-4 sm:mb-5': isAnswer,
                      })}
                      dangerouslySetInnerHTML={{ __html: item }}
                    />
                  );
                })}
              </div>
            </div>
          )}
      </div>
    </div>
  );
}
