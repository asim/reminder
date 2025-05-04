import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber._index';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Fragment } from 'react';
import { Link } from 'react-router';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { ViewMode } from '~/components/quran/view-mode';
import { useQuranViewMode } from '~/hooks/use-quran-view-mode';
import { getChapterOptions } from '~/queries/quran';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(
    getChapterOptions(Number(props.params.chapterNumber))
  );
}

export default function QuranChapter(props: Route.ComponentProps) {
  const { chapterNumber } = props.params;
  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));
  const [mode, setMode] = useQuranViewMode();

  if (!data) {
    return null;
  }

  const previousChapter = Number(chapterNumber) - 1;
  const nextChapter = Number(chapterNumber) + 1;

  return (
    <div className='max-w-4xl flex flex-col w-full mb-8 sm:mb-12 flex-grow mx-auto p-0 lg:p-8'>
      <ViewMode mode={mode} onChange={setMode} />
      <ChapterHeader title={data.name} translation={data.english} subtitle={`Chapter ${data.number}`} />
      {mode === 'arabic' && (
        <div
          dir='rtl'
          className='flex flex-grow flex-wrap font-arabic text-right text-xl sm:text-2xl md:text-3xl leading-loose'
        >
          {data.verses.map((verse) => (
            <Fragment key={verse.number}>
              {verse.arabic.replace('۞', '')}
              &nbsp;۝&nbsp;&nbsp;
            </Fragment>
          ))}
        </div>
      )}
      {mode === 'english' && (
        <div className='flex flex-grow flex-wrap text-left text-base sm:text-lg md:text-xl leading-loose'>
          {data.verses.map((verse) => (
            <Fragment key={verse.number}>{verse.text}&nbsp;</Fragment>
          ))}
        </div>
      )}
      {mode == 'translation' && (
        <div className='space-y-3 sm:space-y-8'>
          {data.verses.map((verse) => (
            <div
              data-chapter-verse={`${data.number}:${verse.number}`}
              key={verse.number}
              className='border-b border-gray-100 pb-3 sm:pb-8'
            >
              <div className='text-xl sm:text-2xl md:text-3xl mb-3 sm:mb-4 text-right leading-loose font-arabic'>
                {verse.arabic.replace('۞', '')}
              </div>
              <div className='text-base sm:text-lg md:text-xl leading-relaxed'>{verse.text}</div>
              <div className='text-xs sm:text-sm text-gray-500 mt-1 sm:mt-2'>
                Verse {verse.number}
              </div>
            </div>
          ))}
        </div>
      )}
      <div className='flex justify-between mt-6 sm:mt-8 mb-3'>
        <PrimaryButton asChild disabled={previousChapter <= 1} className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'>
          <Link to={`/quran/${previousChapter}`}>
            <CircleChevronLeft className='size-4' />
            <span className='ml-1'>Previous</span>
          </Link>
        </PrimaryButton>

        <PrimaryButton asChild disabled={nextChapter >= 114} className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'>
          <Link to={`/quran/${nextChapter}`}>
            <span className='mr-1'>Next</span>
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
