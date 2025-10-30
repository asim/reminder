import { useSuspenseQuery } from '@tanstack/react-query';
import { Link, useParams } from 'react-router';
import { BookmarkButton } from '~/components/interface/bookmark-button';
import { getNameOptions } from '~/queries/names';

export function meta() {
  return [
    { title: 'Names of Allah - Reminder' },
    {
      property: 'og:title',
      content: 'Names of Allah - Reminder',
    },
    {
      name: 'description',
      content:
        'Learn about the Names of Allah, the 99 attributes and qualities through which Muslims identify and connect with Allah',
    },
  ];
}

export default function NameDetail() {
  const params = useParams();
  const nameNumber = parseInt(params.nameNumber || '');

  const { data: name } = useSuspenseQuery(getNameOptions(nameNumber));

  if (!name) {
    return <div>Name not found</div>;
  }

  return (
    <div className='max-w-4xl mx-auto w-full p-0 lg:p-8 mb-8 sm:mb-12 flex-grow'>
      <div className='lg:block hidden mt-0 sm:mt-6 md:mt-8 mb-4 sm:mb-8 md:mb-12'></div>

      <div className='space-y-3 sm:space-y-6'>
        <div className='p-5 sm:p-4 md:p-6 border border-gray-200 rounded-lg space-y-2 sm:space-y-4'>
          <div className='flex flex-col text-left'>
            <div className='flex items-start justify-between'>
              <h1 className='text-xl sm:text-2xl font-semibold text-gray-800'>
                {name.meaning}
              </h1>
              <BookmarkButton
                type='names'
                itemKey={`${name.number}`}
                label={`Name ${name.number}: ${name.meaning}`}
                url={`/names/${name.number}`}
              />
            </div>
            <div className='text-3xl sm:text-4xl md:text-5xl my-4 sm:my-7 font-arabic font-medium text-black'>
              {name.arabic}
            </div>
            <div className='text-lg sm:text-xl md:text-2xl font-medium text-gray-600 mb-1 sm:mb-2'>
              {name.english}
            </div>
          </div>
        </div>

        {name.summary && (
          <div className='p-5 sm:p-4 md:p-6 border border-gray-200 rounded-lg space-y-2 sm:space-y-4'>
            <h2 className='text-lg sm:text-xl font-semibold'>Summary</h2>
            <p className='text-gray-800 text-base sm:text-lg leading-relaxed'>
              {name.summary}
            </p>
          </div>
        )}

        {name.location && name.location.length > 0 && (
          <div className='p-5 sm:p-4 md:p-6 border border-gray-200 rounded-lg space-y-2 sm:space-y-4'>
            <h2 className='text-lg sm:text-xl font-semibold'>
              Locations in Quran
            </h2>
            <div className='flex flex-wrap gap-2'>
              {name.location.filter(Boolean).map((loc, index) => (
                <Link
                  key={index}
                  to={`/quran/${loc.replace(':', '/')}`}
                  className='bg-gray-100 text-sm sm:text-base px-2 sm:px-3 py-1 rounded-full text-gray-700 hover:bg-gray-300 transition-colors'
                >
                  {loc}
                </Link>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
