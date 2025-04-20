import type { Route } from '.react-router/types/app/routes/+types/quran.$chapterNumber';
import { useSuspenseQuery } from '@tanstack/react-query';
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

  if (!data) {
    return null;
  }

  return (
    <div className='max-w-4xl mx-auto'>
      <div className='text-center mb-12'>
        <h1 className='text-4xl font-bold mb-2'>{data.name}</h1>
        <div className='text-xl text-gray-600'>Chapter {data.number}</div>
      </div>

      <div className='space-y-8'>
        {data.verses.map((verse) => (
          <div key={verse.number} className='border-b border-gray-100 pb-8'>
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
    </div>
  );
}
