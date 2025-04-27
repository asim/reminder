import { useSuspenseQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
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
      <div className='text-center mt-8 mb-12'>
        <h1 className='text-4xl font-bold mb-2 flex items-center justify-center'>
          {name.meaning}
        </h1>
        <div className='text-xl text-gray-600 flex flex-col gap-2'>
          <div>{name.english}</div>
          <div className='text-2xl'>{name.arabic}</div>
        </div>
      </div>

      <div className='space-y-8'>
        <div className='p-6 border border-gray-200 rounded-lg shadow-sm space-y-4'>
          <h2 className='text-xl font-semibold'>Description</h2>
          <p className='text-gray-800 text-lg leading-relaxed'>
            {name.description}
          </p>
        </div>

        {name.summary && (
          <div className='p-6 border border-gray-200 rounded-lg shadow-sm space-y-4'>
            <h2 className='text-xl font-semibold'>Summary</h2>
            <p className='text-gray-800 text-lg leading-relaxed'>
              {name.summary}
            </p>
          </div>
        )}

        {name.location && name.location.length > 0 && (
          <div className='p-6 border border-gray-200 rounded-lg shadow-sm space-y-4'>
            <h2 className='text-xl font-semibold'>Found In</h2>
            <div className='flex flex-wrap gap-2'>
              {name.location.map((loc, index) => (
                <span 
                  key={index}
                  className='bg-gray-100 px-3 py-1 rounded-full text-gray-700'
                >
                  {loc}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
