import React from 'react';
import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber.$verseNumber';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Link } from 'react-router';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { ViewMode } from '~/components/quran/view-mode';
import { useQuranViewMode } from '~/hooks/use-quran-view-mode';
import { getChapterOptions } from '~/queries/quran';
import { useWordByWordToggle } from '~/use-word-by-word-toggle';

function toArabicNumber(num: number) {
  return num
    .toString()
    .replace(/\d/g, (d) => String.fromCharCode(0x0660 + Number(d)));
}

export default function QuranVerse(props: Route.ComponentProps) {
  const { chapterNumber, verseNumber } = props.params;

  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));
  const [mode, setMode] = useQuranViewMode();
  const [wordByWord, setWordByWord] = useWordByWordToggle();
  const [showCommentary, setShowCommentary] = React.useState(false);

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
      <ViewMode mode={mode} onChange={setMode} />
      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Verse ${data.number}:${verseNumber}`}
      />
      {mode === 'translation' && (
        <div className='mb-4 flex items-center gap-6'>
          <label className='flex items-center gap-2 cursor-pointer'>
            <input
              type='checkbox'
              checked={wordByWord}
              onChange={e => setWordByWord(e.target.checked)}
              className='accent-black h-4 w-4 rounded'
            />
            <span className='text-sm'>Show word-by-word translation</span>
          </label>
          <label className='flex items-center gap-2 cursor-pointer'>
            <input
              type='checkbox'
              checked={showCommentary}
              onChange={e => setShowCommentary(e.target.checked)}
              className='accent-black h-4 w-4 rounded'
            />
            <span className='text-sm'>Show commentary</span>
          </label>
        </div>
      )}
      
      {mode === 'arabic' && (
        <div className='flex flex-col flex-grow'>
          <div 
            dir='rtl'
            className='text-2xl sm:text-2xl md:text-3xl text-right leading-loose font-arabic'
          >
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
          </div>
        </div>
      )}

      {mode === 'english' && (
        <div className='flex flex-col flex-grow'>
          <div className='text-base sm:text-lg md:text-xl leading-relaxed'>
            {verse.text}
          </div>
        </div>
      )}

      {mode === 'translation' && (
        <div className='flex flex-col flex-grow'>
          <div className='flex flex-row-reverse flex-wrap text-xl sm:text-2xl md:text-3xl mb-3 sm:mb-4 text-right leading-loose font-arabic items-end'>
            {wordByWord && verse.words && verse.words.length > 0
              ? verse.words.map((word, idx, arr) => {
                if (idx === arr.length - 1) {
                  return (
                    <span key={idx} className='verse-arabic-word text-3xl flex flex-col items-center mr-2 mb-2'>
                      <span className='flex flex-row items-center' dir="rtl">
                        <span>{word.arabic}</span>
                        <span className='mx-2 font-arabic'>{toArabicNumber(verse.number)}</span>
                      </span>
                      <span className='text-xs sm:text-sm mt-1 px-1 rounded bg-gray-100 text-gray-700'>{word.english}</span>
                    </span>
                  );
                } else {
                  return (
                    <span key={idx} className='verse-arabic-word text-3xl flex flex-col items-center mr-2 mb-2'>
                      <span>{word.arabic}</span>
                      <span className='text-xs sm:text-sm mt-1 px-1 rounded bg-gray-100 text-gray-700'>{word.english}</span>
                    </span>
                  );
                }
              })
              : (
                <span className='verse-arabic-word text-3xl flex flex-row items-center mr-2 mb-2' dir="rtl">
                  {verse.arabic}
                  <span className='mx-2 font-arabic'>{toArabicNumber(verse.number)}</span>
                </span>
              )}
          </div>
          <div className='text-base sm:text-lg md:text-xl leading-relaxed'>
            {verse.text}
          </div>
          {showCommentary && verse.comments && (
            <div className='mt-2 p-2 bg-yellow-50 border-l-4 border-yellow-400 text-yellow-900 rounded text-sm'>
              <strong>Commentary:</strong> {verse.comments}
            </div>
          )}
        </div>
      )}

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
