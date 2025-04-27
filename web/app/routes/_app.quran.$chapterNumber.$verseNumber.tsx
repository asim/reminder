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
    <div className='max-w-4xl flex flex-col w-full h-full mx-auto'>
      <ChapterHeader
        title={data.name}
        subtitle={`Verse ${data.number}:${verseNumber}`}
      />

      <div className='flex flex-col flex-grow'>
        <div className='text-3xl mb-4 text-right leading-loose font-arabic'>
          {verse.arabic.replace('۞', '')}
        </div>
        <div className='text-xl leading-relaxed'>{verse.text}</div>
      </div>

      <div className='flex justify-between mt-8 mb-12'>
        {hasPreviousVerse ? (
          <PrimaryButton asChild disabled={previousVerse <= 1}>
            <Link to={`/quran/${chapterNumber}/${previousVerse}`}>
              <CircleChevronLeft className='size-4' />
              Previous Verse
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasPreviousChapter}>
            <Link to={`/quran/${previousChapter}`}>
              <CircleChevronLeft className='size-4' />
              Previous Chapter
            </Link>
          </PrimaryButton>
        )}

        {hasNextVerse ? (
          <PrimaryButton asChild disabled={!hasNextVerse}>
            <Link to={`/quran/${chapterNumber}/${nextVerse}`}>
              Next Verse
              <CircleChevronRight className='size-4' />
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasNextChapter}>
            <Link to={`/quran/${nextChapter}`}>
              Next Chapter
              <CircleChevronRight className='size-4' />
            </Link>
          </PrimaryButton>
        )}
      </div>
    </div>
  );
}
