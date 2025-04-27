import type { Route } from '.react-router/types/app/routes/+types/_app.hadith.$book';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
import { getBookOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';
import { useState, useMemo, useEffect } from 'react';
import { Search } from 'lucide-react';

export async function clientLoader(props: Route.LoaderArgs) {
  const bookNumber = parseInt(props.params.book);
  if (isNaN(bookNumber)) return;

  await queryClient.ensureQueryData(getBookOptions(bookNumber));
}

export default function HadithBook() {
  const params = useParams();
  const bookNumber = parseInt(params.book || '');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: book } = useSuspenseQuery(getBookOptions(bookNumber));

  // Reset filter when book changes
  useEffect(() => {
    setSearchQuery('');
  }, [bookNumber]);

  if (!book) {
    return <div>Book not found</div>;
  }

  // Filter hadiths by search query
  const filteredHadiths = useMemo(() => {
    if (!searchQuery) {
      return book.hadiths;
    }

    const query = searchQuery.toLowerCase();
    return book.hadiths.filter(
      (hadith) =>
        hadith.text.toLowerCase().includes(query) ||
        hadith.by.toLowerCase().includes(query) ||
        hadith.info.toLowerCase().includes(query)
    );
  }, [book.hadiths, searchQuery]);

  return (
    <div className='max-w-4xl mx-auto w-full mb-12 flex-grow'>
      <div className='text-center mt-8 mb-12'>
        <h1 className='text-4xl font-bold mb-2 flex items-center justify-center'>
          {book.name}
        </h1>
        <div className='text-xl text-gray-600'>
          Total Ahadith — {book.hadith_count}
        </div>
      </div>

      {/* Search filter */}
      <div className='flex mb-6 justify-start w-full'>
        <div className='relative w-full'>
          <div className='absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none'>
            <Search className='h-4 w-4 text-gray-500' />
          </div>
          <input
            type='text'
            placeholder={`Search in ${book.name}...`}
            className='w-full pl-10 p-2 border border-gray-300 rounded-md'
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button
              onClick={() => setSearchQuery('')}
              className='absolute inset-y-0 right-0 pr-3 flex items-center'
            >
              <span className='text-gray-500 hover:text-gray-700'>×</span>
            </button>
          )}
        </div>
      </div>

      <div className='space-y-4'>
        {filteredHadiths.map((hadith) => (
          <div
            key={`${hadith.info}-${hadith.by}-${bookNumber}`}
            className='p-6 border border-gray-200 rounded-lg space-y-4 hover:border-gray-300 transition-colors'
          >
            <div className='flex justify-between items-center'>
              <span className='font-medium text-gray-700'>
                {hadith.info.trim().replace(/:$/, '')}
              </span>
              <span className='font-medium text-gray-700'>{hadith.by}</span>
            </div>

            <p className='text-gray-800 text-lg leading-relaxed'>
              {hadith.text}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
}
