import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber._index';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Fragment, useEffect } from 'react';
import { Link } from 'react-router';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { ViewMode } from '~/components/quran/view-mode';
import { useQuranViewMode } from '~/hooks/use-quran-view-mode';
import { getChapterOptions } from '~/queries/quran';
import { queryClient } from '~/utils/query-client';

function toArabicNumber(num: number) {
  return num
    .toString()
    .replace(/\d/g, (d) => String.fromCharCode(0x0660 + Number(d)));
}

export function meta() {
  return [
    { title: 'Quran - Reminder' },
    {
      property: 'og:title',
      content: 'Quran - Reminder',
    },
    {
      name: 'description',
      content: 'Read the Quran, the holy book of Islam',
    },
  ];
}

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(
    getChapterOptions(Number(props.params.chapterNumber))
  );
}

export default function QuranChapter(props: Route.ComponentProps) {
  const { chapterNumber } = props.params;
  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));
  const [mode, setMode] = useQuranViewMode();

  useEffect(() => {
    if (!data || !window.location.hash) {
      return;
    }

    const verseId = window.location.hash.substring(1);
    const element = document.getElementById(verseId);
    if (element) {
      setTimeout(() => {
        element.scrollIntoView({ behavior: 'instant' });
      }, 100);
    }
  }, [data, mode]);

  if (!data) {
    return null;
  }

  const previousChapter = Number(chapterNumber) - 1;
  const nextChapter = Number(chapterNumber) + 1;

  return (
    <div className='max-w-4xl flex flex-col w-full mb-8 sm:mb-12 flex-grow mx-auto p-0 lg:p-8'>
      <ViewMode mode={mode} onChange={setMode} />
      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Chapter ${data.number}`}
      />
      {mode === 'arabic' && (
        <div
          dir='rtl'
          className='flex flex-grow flex-wrap font-arabic text-right text-2xl sm:text-2xl md:text-3xl leading-loose content-start'
        >
          {data.verses.map((verse) => (
            <Fragment key={verse.number}>
              {verse.words && verse.words.length > 0
                ? verse.words.map((word, idx, arr) => (
                    <span key={idx} className='verse-arabic-word'>
                      {word.arabic}
                      {idx === arr.length - 1 && (
                        <span className='mx-2 font-arabic'>
                          {toArabicNumber(verse.number)}
                        </span>
                      )}
                      &nbsp;
                    </span>
                  ))
                : verse.arabic.split(' ').map((word, idx, arr) => (
                    <span key={idx} className='verse-arabic-word'>
                      {word}
                      {idx === arr.length - 1 && (
                        <span className='mx-2 font-arabic'>
                          {toArabicNumber(verse.number)}
                        </span>
                      )}
                      &nbsp;
                    </span>
                  ))}
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
              id={`${verse.number}`}
              data-chapter-verse={`${data.number}:${verse.number}`}
              key={verse.number}
              className='border-b border-gray-100 pb-3 sm:pb-8'
            >
              <div className='flex flex-row-reverse flex-wrap text-xl sm:text-2xl md:text-3xl mb-3 sm:mb-4 text-right leading-loose font-arabic items-end'>
                {verse.words && verse.words.length > 0
                  ? verse.words.map((word, idx, arr) => {
                      if (idx === arr.length - 1) {
                        return (
                          <span key={idx} className='verse-arabic-word text-2xl flex flex-col items-center mr-2 mb-2'>
                            <span className='flex flex-row items-center' dir="rtl">
                              <span>{word.arabic}</span>
                              <span className='mx-2 font-arabic'>{toArabicNumber(verse.number)}</span>
                            </span>
                            <span className='text-xs sm:text-sm mt-1 px-1 rounded bg-gray-100 text-gray-700'>{word.english}</span>
                          </span>
                        );
                      } else {
                        return (
                          <span key={idx} className='verse-arabic-word text-2xl flex flex-col items-center mr-2 mb-2'>
                            <span>{word.arabic}</span>
                            <span className='text-xs sm:text-sm mt-1 px-1 rounded bg-gray-100 text-gray-700'>{word.english}</span>
                          </span>
                        );
                      }
                    })
                  : verse.arabic.split(' ').map((word, idx, arr) => (
                      idx === arr.length - 1 ? (
                        <span key={idx} className='verse-arabic-word text-2xl flex flex-row items-center mr-2 mb-2' dir="rtl">
                          <span>{word}</span>
                          <span className='mx-2 font-arabic'>{toArabicNumber(verse.number)}</span>
                        </span>
                      ) : (
                        <span key={idx} className='verse-arabic-word text-2xl mr-2 mb-2'>
                          {word}
                        </span>
                      )
                    ))
                }
              </div>
              <div className='text-base sm:text-lg md:text-xl leading-relaxed'>
                {verse.text}
              </div>
              <a
                href={verse.number !== 0 ? `/quran/${data.number}/${verse.number}` : undefined}
                className={`text-xs sm:text-sm text-gray-500 mt-1 sm:mt-2${verse.number === 0 ? ' hidden' : ''}`}
                style={verse.number === 0 ? { display: 'none' } : {}}
              >
                {verse.number !== 0 ? `${data.number}:${verse.number}` : ''}
              </a>
            </div>
          ))}
        </div>
      )}
      <div className='flex justify-between mt-6 sm:mt-8 mb-3'>
        <PrimaryButton
          asChild
          disabled={previousChapter <= 1}
          className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'
        >
          <Link to={`/quran/${previousChapter}`}>
            <CircleChevronLeft className='size-4' />
            <span className='ml-1'>Previous</span>
          </Link>
        </PrimaryButton>

        <PrimaryButton
          asChild
          disabled={nextChapter >= 114}
          className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'
        >
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
