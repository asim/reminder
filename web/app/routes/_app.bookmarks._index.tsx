import { Trash2 } from 'lucide-react';
import { useEffect } from 'react';
import { Link } from 'react-router';
import { useBookmarks } from '~/hooks/use-bookmarks';
import { httpGet } from '~/utils/http';

export function meta() {
  return [
    { title: 'Bookmarks - Reminder' },
    {
      property: 'og:title',
      content: 'Bookmarks - Reminder',
    },
    {
      name: 'description',
      content: 'Your saved verses, hadiths, and names of Allah',
    },
  ];
}

export default function BookmarksPage() {
  const { bookmarks, removeBookmark, updateExcerpt } = useBookmarks();

  // Migrate existing bookmarks by fetching excerpts
  useEffect(() => {
    const migrateQuranExcerpts = async () => {
      for (const [key, bookmark] of Object.entries(bookmarks.quran)) {
        if (bookmark.excerpt) continue; // Already has excerpt
        
        const [chapter, verse] = key.split(':').map(Number);
        if (!chapter || !verse) continue;
        
        try {
          const data = await httpGet<{ verses: { number: number; text: string }[] }>(`/api/quran/${chapter}`);
          const verseData = data.verses.find(v => v.number === verse);
          if (verseData) {
            const excerpt = verseData.text.length > 80 ? verseData.text.slice(0, 80) + '...' : verseData.text;
            updateExcerpt('quran', key, excerpt);
          }
        } catch (e) {
          console.error('Failed to fetch excerpt for', key, e);
        }
      }
    };

    const migrateHadithExcerpts = async () => {
      for (const [key, bookmark] of Object.entries(bookmarks.hadith)) {
        if (bookmark.excerpt) continue;
        
        const [book, hadithNum] = key.split(':').map(Number);
        if (!book || !hadithNum) continue;
        
        try {
          const data = await httpGet<{ hadiths: { number: number; text: string }[] }>(`/api/hadith/${book}`);
          const hadith = data.hadiths.find(h => h.number === hadithNum);
          if (hadith) {
            const excerpt = hadith.text.length > 80 ? hadith.text.slice(0, 80) + '...' : hadith.text;
            updateExcerpt('hadith', key, excerpt);
          }
        } catch (e) {
          console.error('Failed to fetch excerpt for', key, e);
        }
      }
    };

    migrateQuranExcerpts();
    migrateHadithExcerpts();
  }, []); // Run once on mount

  const quranEntries = Object.entries(bookmarks.quran);
  const hadithEntries = Object.entries(bookmarks.hadith);
  const namesEntries = Object.entries(bookmarks.names);

  const hasAnyBookmarks =
    quranEntries.length > 0 ||
    hadithEntries.length > 0 ||
    namesEntries.length > 0;

  return (
    <div className='max-w-4xl mx-auto w-full p-4 lg:p-8 mb-8 sm:mb-12 flex-grow overflow-y-auto'>
      <div className='text-center mt-0 sm:mt-6 md:mt-8 mb-4 sm:mb-8 md:mb-12'>
        <h1 className='text-2xl sm:text-3xl md:text-4xl font-bold mb-1 sm:mb-2'>
          Bookmarks
        </h1>
        <div className='text-base sm:text-lg md:text-xl text-gray-600'>
          Your saved verses, hadiths, and names of Allah
        </div>
      </div>

      {!hasAnyBookmarks && (
        <div className='text-center py-12 text-gray-500'>
          <p className='text-lg mb-2'>No bookmarks yet</p>
          <p className='text-sm'>
            Start saving verses, hadiths, and names by clicking the star icon
          </p>
        </div>
      )}

      {quranEntries.length > 0 && (
        <div className='mb-8'>
          <h2 className='text-xl sm:text-2xl font-semibold mb-4'>Quran</h2>
          <div className='space-y-2'>
            {quranEntries.map(([key, bookmark]) => (
              <div
                key={key}
                className='flex items-center justify-between p-3 sm:p-4 border border-gray-200 rounded-lg hover:border-gray-300 transition-colors'
              >
                <Link
                  to={bookmark.url}
                  className='flex-grow text-gray-800 hover:text-black'
                >
                  <div className='text-base sm:text-lg font-medium'>{bookmark.label}</div>
                  {bookmark.excerpt && (
                    <div className='text-sm text-gray-500 mt-1 line-clamp-2'>{bookmark.excerpt}</div>
                  )}
                </Link>
                <button
                  onClick={() => removeBookmark('quran', key)}
                  className='ml-4 p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors flex-shrink-0'
                  aria-label='Delete bookmark'
                >
                  <Trash2 className='size-4' />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {hadithEntries.length > 0 && (
        <div className='mb-8'>
          <h2 className='text-xl sm:text-2xl font-semibold mb-4'>Hadith</h2>
          <div className='space-y-2'>
            {hadithEntries.map(([key, bookmark]) => (
              <div
                key={key}
                className='flex items-center justify-between p-3 sm:p-4 border border-gray-200 rounded-lg hover:border-gray-300 transition-colors'
              >
                <Link
                  to={bookmark.url}
                  className='flex-grow text-gray-800 hover:text-black'
                >
                  <div className='text-base sm:text-lg font-medium'>{bookmark.label}</div>
                  {bookmark.excerpt && (
                    <div className='text-sm text-gray-500 mt-1 line-clamp-2'>{bookmark.excerpt}</div>
                  )}
                </Link>
                <button
                  onClick={() => removeBookmark('hadith', key)}
                  className='ml-4 p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors flex-shrink-0'
                  aria-label='Delete bookmark'
                >
                  <Trash2 className='size-4' />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}

      {namesEntries.length > 0 && (
        <div className='mb-8'>
          <h2 className='text-xl sm:text-2xl font-semibold mb-4'>
            Names of Allah
          </h2>
          <div className='space-y-2'>
            {namesEntries.map(([key, bookmark]) => (
              <div
                key={key}
                className='flex items-center justify-between p-3 sm:p-4 border border-gray-200 rounded-lg hover:border-gray-300 transition-colors'
              >
                <Link
                  to={bookmark.url}
                  className='flex-grow text-gray-800 hover:text-black'
                >
                  <div className='text-base sm:text-lg font-medium'>{bookmark.label}</div>
                  {bookmark.excerpt && (
                    <div className='text-sm text-gray-500 mt-1 line-clamp-2'>{bookmark.excerpt}</div>
                  )}
                </Link>
                <button
                  onClick={() => removeBookmark('names', key)}
                  className='ml-4 p-2 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-md transition-colors flex-shrink-0'
                  aria-label='Delete bookmark'
                >
                  <Trash2 className='size-4' />
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
