import { useSuspenseQuery } from '@tanstack/react-query';
import { listBooksOptions } from '~/queries/hadith';

export default function HadithIndex() {
  const { data: books } = useSuspenseQuery(listBooksOptions());

  return (
    <div className='flex flex-col flex-1 items-center justify-center'>
      <h1 className='text-2xl font-normal text-gray-400'>
        Please select a book
      </h1>
    </div>
  );
}
