import type { Route } from '.react-router/types/app/routes/+types/quran.$chapterNumber._index';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Fragment, useState } from 'react';
import { Link } from 'react-router';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { getChapterOptions } from '~/queries/quran';
import { cn } from '~/utils/classname';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(
    getChapterOptions(Number(props.params.chapterNumber))
  );
}

export default function QuranChapter(props: Route.ComponentProps) {
  const { chapterNumber } = props.params;
  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));

  const [mode, setMode] = useState<'translation' | 'recitation'>('translation');

  if (!data) {
    return null;
  }

  const previousChapter = Number(chapterNumber) - 1;
  const nextChapter = Number(chapterNumber) + 1;

  return (
    <div className='max-w-4xl flex flex-col w-full h-full mx-auto'>
      <div className='flex items-center justify-center gap-2 mb-8'>
        <button
          type='button'
          onClick={() => setMode('translation')}
          className={cn(
            'px-4 py-2 rounded-md cursor-pointer',
            mode === 'translation' && 'bg-gray-100'
          )}
        >
          Translation Mode
        </button>
        <button
          type='button'
          onClick={() => setMode('recitation')}
          className={cn(
            'px-4 py-2 rounded-md cursor-pointer',
            mode === 'recitation' && 'bg-gray-100'
          )}
        >
          Recitation Mode
        </button>
      </div>

      <ChapterHeader title={data.name} subtitle={`Chapter ${data.number}`} />

      {mode === 'recitation' && (
        <div
          dir='rtl'
          className='flex flex-grow flex-wrap font-arabic text-right text-3xl leading-loose'
        >
          {data.verses.map((verse) => (
            <Fragment key={verse.number}>
              {verse.arabic}
              &nbsp;۝&nbsp;&nbsp;
            </Fragment>
          ))}
        </div>
      )}

      {mode == 'translation' && (
        <div className='space-y-8'>
          {data.verses.map((verse) => (
            <div
              data-chapter-verse={`${data.number}:${verse.number}`}
              key={verse.number}
              className='border-b border-gray-100 pb-8'
            >
              <div className='text-3xl mb-4 text-right leading-loose font-arabic'>
                {verse.arabic}
              </div>
              <div className='text-xl leading-relaxed'>{verse.text}</div>
              <div className='text-sm text-gray-500 mt-2'>
                Verse {verse.number}
              </div>
            </div>
          ))}
        </div>
      )}

      <div className='flex justify-between mt-8 mb-12'>
        <PrimaryButton asChild disabled={previousChapter <= 1}>
          <Link to={`/quran/${previousChapter}`}>
            <CircleChevronLeft className='size-4' />
            Previous Chapter
          </Link>
        </PrimaryButton>

        <PrimaryButton asChild disabled={nextChapter >= 114}>
          <Link to={`/quran/${nextChapter}`}>
            Next Chapter
            <CircleChevronRight className='size-4' />
          </Link>
        </PrimaryButton>
      </div>
    </div>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  return <PageError error={error} />;
}
