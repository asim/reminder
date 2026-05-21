import { useSuspenseQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
import { getChapterOptions } from '~/queries/quran';

function toArabicNumber(num: number) {
  return num
    .toString()
    .replace(/\d/g, (d) => String.fromCharCode(0x0660 + Number(d)));
}

export default function EmbedQuranVerse() {
  const { chapter, verse: verseNum } = useParams();
  const chapterNumber = Number(chapter);
  const verseNumber = Number(verseNum);

  const { data } = useSuspenseQuery(getChapterOptions(chapterNumber));

  const verse = data?.verses.find((v) => v.number === verseNumber);

  if (!verse) {
    return <div className='p-4 text-sm text-gray-500'>Verse not found</div>;
  }

  return (
    <div className='font-sans p-4 max-w-lg'>
      {verse.arabic && (
        <div
          dir='rtl'
          className='text-xl leading-loose font-arabic text-right text-gray-800 mb-3'
        >
          {verse.arabic}
          <span className='mx-1 font-arabic'>﴿{toArabicNumber(verseNumber)}﴾</span>
        </div>
      )}
      <p className='text-sm text-gray-700 leading-relaxed mb-3'>
        {verse.text}
      </p>
      <div className='flex items-center justify-between border-t border-gray-100 pt-2'>
        <a
          href={`https://reminder.dev/quran/${chapterNumber}/${verseNumber}`}
          target='_blank'
          rel='noopener noreferrer'
          className='text-xs text-gray-400 hover:text-gray-600 no-underline'
        >
          {data.english} {chapterNumber}:{verseNumber} — reminder.dev
        </a>
      </div>
    </div>
  );
}
