import { useSuspenseQuery } from '@tanstack/react-query';
import { Link, useParams } from 'react-router';
import { getNameOptions } from '~/queries/names';

export default function NameDetail() {
  const params = useParams();
  const nameNumber = parseInt(params.nameNumber || '');

  const { data: name } = useSuspenseQuery(getNameOptions(nameNumber));

  if (!name) {
    return <div>Name not found</div>;
  }

  return (
    <div className='max-w-4xl mx-auto w-full'>
      <div className='mt-8 mb-12'></div>

      <div className='space-y-6'>
        <div className='p-6 border border-gray-200 rounded-lg space-y-4'>
          <div className='flex flex-col text-left'>
            <h1 className='text-2xl font-semibold text-gray-800'>
              {name.meaning}
            </h1>
            <div className='text-5xl my-7 font-arabic font-medium text-black'>
              {name.arabic}
            </div>
            <div className='text-2xl font-medium text-gray-600 mb-2'>
              {name.english}
            </div>
          </div>
          <p className='text-gray-800 text-lg leading-relaxed'>
            {name.description.replace(/\\"/g, '')}
          </p>
        </div>

        {name.summary && (
          <div className='p-6 border border-gray-200 rounded-lg space-y-4'>
            <h2 className='text-xl font-semibold'>Summary</h2>
            <p className='text-gray-800 text-lg leading-relaxed'>
              {name.summary}
            </p>
          </div>
        )}

        {name.location && name.location.length > 0 && (
          <div className='p-6 border border-gray-200 rounded-lg space-y-4'>
            <h2 className='text-xl font-semibold'>Locations in Quran</h2>
            <div className='flex flex-wrap gap-2'>
              {name.location.filter(Boolean).map((loc, index) => (
                <Link
                  key={index}
                  to={`/quran/${loc.replace(':', '/')}`}
                  className='bg-gray-100 px-3 py-1 rounded-full text-gray-700 hover:bg-gray-300 transition-colors'
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
