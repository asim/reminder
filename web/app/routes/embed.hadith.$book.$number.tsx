import { useSuspenseQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
import { getBookOptions } from '~/queries/hadith';

export default function EmbedHadith() {
  const { book, number } = useParams();
  const bookNumber = Number(book);
  const hadithNumber = Number(number);

  const { data } = useSuspenseQuery(getBookOptions(bookNumber));

  const hadith = data?.hadiths.find((h) => h.number === hadithNumber);

  if (!hadith) {
    return <div className='p-4 text-sm text-gray-500'>Hadith not found</div>;
  }

  const narrator = hadith.narrator || hadith.by || '';
  const text = hadith.english || hadith.text || '';

  return (
    <div className='font-sans p-4 max-w-lg'>
      {narrator && (
        <p className='text-xs text-gray-500 mb-2'>{narrator}</p>
      )}
      <p className='text-sm text-gray-700 leading-relaxed mb-3'>
        {text}
      </p>
      <div className='border-t border-gray-100 pt-2'>
        <a
          href={`https://reminder.dev/hadith/${bookNumber}#${hadithNumber}`}
          target='_blank'
          rel='noopener noreferrer'
          className='text-xs text-gray-400 hover:text-gray-600 no-underline'
        >
          {data.name}, Hadith {hadithNumber} — reminder.dev
        </a>
      </div>
    </div>
  );
}
