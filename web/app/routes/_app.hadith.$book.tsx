import type { Route } from '.react-router/types/app/routes/+types/_app.hadith.$book';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
import { getBookOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  const bookNumber = parseInt(props.params.book);
  if (isNaN(bookNumber)) return;

  await queryClient.ensureQueryData(getBookOptions(bookNumber));
}

export default function HadithBook() {
  const params = useParams();
  const bookNumber = parseInt(params.book || '');

  const { data: book } = useSuspenseQuery(getBookOptions(bookNumber));

  if (!book) {
    return <div>Book not found</div>;
  }

  return (
    <div className='max-w-4xl mx-auto'>
      <div className='text-center mt-8 mb-12'>
        <h1 className='text-4xl font-bold mb-2 flex items-center justify-center'>
          {book.name}
        </h1>
        <div className='text-xl text-gray-600'>
          Total Ahadith — {book.hadith_count}
        </div>
      </div>

      <div className='space-y-8'>
        {book.hadiths.map((hadith) => (
          <div
            key={`${hadith.info}-${hadith.by}-${bookNumber}`}
            className='p-6 border border-gray-200 rounded-lg shadow-sm space-y-4 hover:border-gray-300 transition-colors'
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
