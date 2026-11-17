import { ExplainButton } from '~/components/quran/explain-button';
import { collectionFromPath, collectionTitle, hadithPath } from '~/utils/hadith-collection';
import { readingMeta } from '~/utils/reading-meta';
import { readingLoader } from '~/utils/prerender.server';
import type { Route } from '.react-router/types/app/routes/+types/_app.hadith.$book';
import { useSuspenseQuery } from '@tanstack/react-query';
import { Search } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { useLocation, useParams } from 'react-router';
import { BookmarkButton } from '~/components/interface/bookmark-button';
import { ReadingBookmarkButton } from '~/components/interface/reading-bookmark-button';
import { ShareButton } from '~/components/interface/share-button';
import { getBookOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

function toArabicNumber(num: number) {
  return num
    .toString()
    .replace(/\d/g, (d) => String.fromCharCode(0x0660 + Number(d)));
}



export async function clientLoader(props: Route.LoaderArgs) {
  const bookNumber = parseInt(props.params.book);
  if (isNaN(bookNumber)) return;

  const options = getBookOptions(bookNumber, collectionFromPath(new URL(props.request.url).pathname));
  const data = await queryClient.ensureQueryData(options);
  return { readingQueries: [{ queryKey: options.queryKey, data }] };
}

export default function HadithBook() {
  const params = useParams();
  const collection = collectionFromPath(useLocation().pathname);
  const prefix = hadithPath(collection);
  const title = collectionTitle(collection);
  const recordLabel = collection === 'muslim' ? 'Record' : 'Hadith';
  const bookNumber = parseInt(params.book || '');
  const [searchQuery, setSearchQuery] = useState('');

  const { data: book } = useSuspenseQuery(getBookOptions(bookNumber, collection));

  // Reset filter when book changes
  useEffect(() => {
    setSearchQuery('');
  }, [bookNumber, collection]);

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
        (hadith.english || hadith.text || '').toLowerCase().includes(query) ||
        (hadith.narrator || hadith.by || '').toLowerCase().includes(query) ||
        (hadith.arabic || '').includes(query)
    );
  }, [book.hadiths, searchQuery]);

  return (
    <div className='max-w-4xl mx-auto w-full mb-8 sm:mb-12 flex-grow p-0 lg:p-8'>
      <div className='flex justify-end mt-0 sm:mt-6 md:mt-8'>
        <ShareButton
          title={`${book.name} - ${title}`}
          url={`${prefix}/${bookNumber}`}
        />
      </div>
      <div className='text-center mb-4 sm:mb-8 md:mb-12'>
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-bold mb-1 sm:mb-2 flex items-center justify-center'>
          {book.name}
        </h1>
        <div className='text-base sm:text-lg md:text-xl text-gray-600'>
          {title} · {book.hadith_count} narrations
        </div>
      </div>

      {collection === 'muslim' && <p className='mb-4 text-xs text-gray-500'>Source: <a className='underline' href='https://github.com/asim/hadith-json'>hadith-json</a>. Record numbers follow this dataset.</p>}
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
          const hadithNumber = hadith.number || idx + 1;
          const narrator = hadith.narrator || hadith.by || '';
          const englishText = hadith.english || hadith.text || '';
          const arabicText = hadith.arabic || '';
          const excerptText = englishText || arabicText;

          return (
            <div
              key={`${collection}-${bookNumber}-${hadithNumber}`}
              id={String(hadithNumber)}
              data-hadith={`${collection}:${bookNumber}:${hadithNumber}`}
              className='border-b border-gray-100 pb-3 sm:pb-8 space-y-2 sm:space-y-4 scroll-mt-4'
            >
              {narrator && <p className='text-sm text-gray-500'>{narrator}</p>}

              {/* Arabic text */}
              {arabicText && (
                <div
                  lang='ar' dir='rtl'
                  className='text-xl sm:text-2xl md:text-3xl leading-loose font-arabic text-right text-gray-800 border-b border-gray-100 pb-4'
                >
                  {arabicText}
                  <span className='mx-2 font-arabic'>
                    ﴿{toArabicNumber(hadithNumber)}﴾
                  </span>
                </div>
              )}

              {/* English translation */}
              <p className='text-gray-800 text-base sm:text-lg leading-relaxed'>
                {englishText || 'English translation is unavailable in this source.'}
              </p>
              <div className='flex items-center justify-between mt-2 gap-2'>
                <a className='text-xs sm:text-sm text-gray-500' href={`${prefix}/${bookNumber}#${hadithNumber}`}>{recordLabel} {hadithNumber}</a>
                <div className='flex items-center gap-1'>
                  <ExplainButton hadith={{ collection, book: bookNumber, number: hadithNumber, label: `${title} · ${book.name} — ${recordLabel} ${hadithNumber}` }} />
                  <ReadingBookmarkButton
                    type='hadith'
                    label={`${title} · ${book.name} - ${recordLabel} ${hadithNumber}`}
                    url={`${prefix}/${bookNumber}#${hadithNumber}`}
                    excerpt={excerptText.slice(0, 80)}
                  />
                  <BookmarkButton
                    type='hadith'
                    itemKey={collection === 'bukhari' ? `${bookNumber}:${hadithNumber}` : `muslim:${bookNumber}:${hadithNumber}`}
                    label={`${title} · ${book.name} - ${recordLabel} ${hadithNumber}`}
                    url={`${prefix}/${bookNumber}#${hadithNumber}`}
                    excerpt={excerptText.slice(0, 80)}
                  />
                  <ShareButton
                    title={`${title} · ${book.name} - ${recordLabel} ${hadithNumber}`}
                    text={englishText}
                    url={`${prefix}/${bookNumber}#${hadithNumber}`}
                  />
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}

export function loader({ params, request }: { params: Record<string, string | undefined>; request: Request }) {
  const collection = collectionFromPath(new URL(request.url).pathname);
  return readingLoader(collection === 'muslim' ? 'muslim' : 'hadith', ['get-hadith-book', collection, Number(params.book)], params.book);
}

export const meta = (args: any) => readingMeta('hadith', args);
