import type { Route } from '.react-router/types/app/routes/+types/_app.hadith.$book';
import { useSuspenseQuery } from '@tanstack/react-query';
import { Search } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { useParams } from 'react-router';
import { BookmarkButton } from '~/components/interface/bookmark-button';
import { getBookOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

export function meta() {
  return [
    { title: 'Hadith - Reminder' },
    {
      property: 'og:title',
      content: 'Hadith - Reminder',
    },
    {
      name: 'description',
      content:
        'Read the Hadith, the sayings of Prophet Muhammad (Peace Be Upon Him)',
    },
  ];
}

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

  // Scroll to hash if present
  useEffect(() => {
    if (!book || !window.location.hash) {
      return;
    }

    const hadithId = window.location.hash.substring(1);
    const element = document.getElementById(hadithId);
    if (element) {
      setTimeout(() => {
        // Find the scrollable parent container
        const scrollContainer = element.closest('.overflow-y-auto');
        if (scrollContainer) {
          // Get the element's position relative to the scroll container
          const containerRect = scrollContainer.getBoundingClientRect();
          const elementRect = element.getBoundingClientRect();
          const relativeTop = elementRect.top - containerRect.top;
          const currentScrollTop = scrollContainer.scrollTop;
          const targetScrollTop = currentScrollTop + relativeTop - 100;

          scrollContainer.scrollTo({ top: targetScrollTop, behavior: 'smooth' });
        } else {
          element.scrollIntoView({ behavior: 'smooth', block: 'start' });
        }
      }, 100);
    }
  }, [book]);

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
    <div className='max-w-4xl mx-auto w-full mb-8 sm:mb-12 flex-grow p-0 lg:p-8'>
      <div className='text-center mt-0 sm:mt-6 md:mt-8 mb-4 sm:mb-8 md:mb-12'>
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-bold mb-1 sm:mb-2 flex items-center justify-center'>
          {book.name}
        </h1>
        <div className='text-base sm:text-lg md:text-xl text-gray-600'>
          Total Ahadith — {book.hadith_count}
        </div>
      </div>

      {/* Search filter */}
      <div className='flex mb-4 sm:mb-6 justify-start w-full'>
        <div className='relative w-full'>
          <div className='absolute inset-y-0 left-0 pl-2 sm:pl-3 flex items-center pointer-events-none'>
            <Search className='h-3 w-3 sm:h-4 sm:w-4 text-gray-500' />
          </div>
          <input
            type='text'
            placeholder={`Search in ${book.name}...`}
            className='w-full text-sm sm:text-base pl-8 sm:pl-10 p-1.5 sm:p-2 border border-gray-300 rounded-md'
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && (
            <button
              onClick={() => setSearchQuery('')}
              className='absolute inset-y-0 right-0 pr-2 sm:pr-3 flex items-center'
            >
              <span className='text-gray-500 hover:text-gray-700'>×</span>
            </button>
          )}
        </div>
      </div>

      <div className='space-y-3 sm:space-y-4'>
        {filteredHadiths.map((hadith, idx) => {
          // Extract hadith number from info field (e.g., "Volume 9, Book 93, Number 473 :")
          const hadithNumber = hadith.info.includes('Number')
            ? hadith.info.split('Number')[1].split(':')[0].trim()
            : `${idx + 1}`;

          return (
            <div
              key={`${hadith.info}-${hadith.by}-${bookNumber}`}
              id={hadithNumber}
              className='p-3 sm:p-4 md:p-6 border border-gray-200 rounded-lg space-y-2 sm:space-y-4 hover:border-gray-300 transition-colors'
            >
              <div className='flex gap-2 lg:mb-0 mb-4 lg:gap-0 lg:flex-row flex-col justify-start lg:justify-between items-start lg:items-center'>
                <span className='text-sm sm:text-base font-medium text-gray-700'>
                  {hadith.info.trim().replace(/:$/, '')}
                </span>
                <div className='flex items-center gap-2'>
                  <span className='text-sm text-balance sm:text-base font-medium text-gray-700'>
                    {hadith.by}
                  </span>
                  <BookmarkButton
                    type='hadith'
                    itemKey={`${bookNumber}:${hadithNumber}`}
                    label={`Hadith ${bookNumber}:${hadithNumber} - ${hadith.info}`}
                    url={`/hadith/${bookNumber}#${hadithNumber}`}
                    excerpt={hadith.text.length > 80 ? hadith.text.slice(0, 80) + '...' : hadith.text}
                  />
                </div>
              </div>

              <p className='text-gray-800 text-base sm:text-lg leading-relaxed'>
                {hadith.text}
              </p>
            </div>
          );
        })}
      </div>
    </div>
  );
}
