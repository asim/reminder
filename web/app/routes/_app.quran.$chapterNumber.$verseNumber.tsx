import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber.$verseNumber';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Link } from 'react-router';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { getChapterOptions } from '~/queries/quran';

export default function QuranVerse(props: Route.ComponentProps) {
  const { chapterNumber, verseNumber } = props.params;

  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));

  const verse = data.verses.find(
    (verse) => verse.number === Number(verseNumber)
  );

  if (!verse) {
    return <PageError error={new Error('Verse not found')} />;
  }

  const previousVerse = Number(verseNumber) - 1;
  const nextVerse = Number(verseNumber) + 1;
  const hasPreviousVerse = previousVerse >= 1;
  const hasNextVerse = nextVerse <= data.verses.length;

  const nextChapter = Number(chapterNumber) + 1;
  const previousChapter = Number(chapterNumber) - 1;
  const hasNextChapter = nextChapter <= 114;
  const hasPreviousChapter = previousChapter >= 1;

  return (
    <div className='max-w-4xl flex flex-col w-full mx-auto p-0 lg:p-8'>
      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Verse ${data.number}:${verseNumber}`}
      />

      <div className='flex flex-col flex-grow'>
        <div className='text-2xl sm:text-2xl md:text-3xl mb-3 sm:mb-4 text-right leading-loose font-arabic'>
          {verse.arabic.replace('۞', '')}
        </div>
        <div className='text-base sm:text-lg md:text-xl leading-relaxed'>{verse.text}</div>
      </div>

      <div className='flex justify-between mt-6 sm:mt-8 mb-8 sm:mb-12'>
        {hasPreviousVerse ? (
          <PrimaryButton asChild disabled={previousVerse <= 1} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={`/quran/${chapterNumber}/${previousVerse}`}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasPreviousChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={`/quran/${previousChapter}`}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        )}

        {hasNextVerse ? (
          <PrimaryButton asChild disabled={!hasNextVerse} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={`/quran/${chapterNumber}/${nextVerse}`}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasNextChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={`/quran/${nextChapter}`}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        )}
      </div>
    </div>
  );
}
