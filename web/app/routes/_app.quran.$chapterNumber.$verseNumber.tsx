import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber.$verseNumber';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import React from 'react';
import { Link, useSearchParams } from 'react-router';
import { EmbedButton } from '~/components/interface/embed-button';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ShareButton } from '~/components/interface/share-button';
import { AudioPlayer } from '~/components/quran/audio-player';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { VerseArabic } from '~/components/quran/verse-arabic';
import { ViewMode } from '~/components/quran/view-mode';
import { useQuranViewMode } from '~/hooks/use-quran-view-mode';
import { getChapterOptions } from '~/queries/quran';
import { useTransliterationToggle } from '~/use-transliteration-toggle';
import { useWordByWordToggle } from '~/use-word-by-word-toggle';
import { buildQuranShareUrl } from '~/utils/quran-share';


export default function QuranVerse(props: Route.ComponentProps) {
  const { chapterNumber, verseNumber } = props.params;

  const { data } = useSuspenseQuery(getChapterOptions(Number(chapterNumber)));
  const [searchParams] = useSearchParams();
  const [mode, setMode] = useQuranViewMode();
  const [wordByWord, setWordByWord] = useWordByWordToggle();
  const [showTransliteration, setShowTransliteration] = useTransliterationToggle();
  const [showCommentary, setShowCommentary] = React.useState(() => {
    const param = searchParams.get('commentary');
    return param === '1' || param === 'true';
  });

  const verse = data.verses.find(
    (verse) => verse.number === Number(verseNumber)
  );

  if (!verse) {
    return <PageError error={new Error('Verse not found')} />;
  }

  const previousVerse = Number(verseNumber) - 1;
  const nextVerse = Number(verseNumber) + 1;
  const hasPreviousVerse = previousVerse >= 1;

  // Chapters carry the Bismillah as an entry numbered 0, so the array is one
  // longer than the verse count. Al-Fatiha is exempt because its first verse
  // is the Bismillah, and At-Tawbah because it opens without one. Use the
  // highest verse number, or the last verse links past the end of the chapter.
  const lastVerse = data.verses.reduce(
    (max, v) => (v.number > max ? v.number : max),
    0
  );
  const hasNextVerse = nextVerse <= lastVerse;

  const nextChapter = Number(chapterNumber) + 1;
  const previousChapter = Number(chapterNumber) - 1;
  const hasNextChapter = nextChapter <= 114;
  const hasPreviousChapter = previousChapter >= 1;

  return (
    <div className='max-w-4xl flex flex-col w-full mx-auto p-0 lg:p-8'>
      <ViewMode mode={mode} onChange={setMode} />
      <div className='flex justify-end gap-1'>
        <EmbedButton path={`/quran/${data.number}/${verseNumber}`} />
        <ShareButton
          title={`Quran ${data.number}:${verseNumber} - ${data.english}`}
          text={verse.text}
          url={buildQuranShareUrl(`/quran/${data.number}/${verseNumber}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}
        />
      </div>
      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Verse ${data.number}:${verseNumber}`}
      />



      {(mode === 'translation' || mode === 'arabic') && (
        <div className='mb-4 flex items-center gap-4 sm:gap-6 flex-wrap'>
          {mode === 'translation' && (
            <label className='flex items-center gap-2 cursor-pointer'>
              <input
                type='checkbox'
                checked={wordByWord}
                onChange={e => setWordByWord(e.target.checked)}
                className='accent-black h-4 w-4 rounded'
              />
              <span className='text-sm'>Word-by-word</span>
            </label>
          )}
          <label className='flex items-center gap-2 cursor-pointer'>
            <input
              type='checkbox'
              checked={showTransliteration}
              onChange={e => setShowTransliteration(e.target.checked)}
              className='accent-black h-4 w-4 rounded'
            />
            <span className='text-sm'>Transliteration</span>
          </label>
          {mode === 'translation' && (
            <label className='flex items-center gap-2 cursor-pointer'>
              <input
                type='checkbox'
                checked={showCommentary}
                onChange={e => setShowCommentary(e.target.checked)}
                className='accent-black h-4 w-4 rounded'
              />
              <span className='text-sm'>Commentary</span>
            </label>
          )}
        </div>
      )}
      {/* Audio Player */}
      <AudioPlayer
        arabicUrl={verse.audio_arabic}
        englishUrl={verse.audio_english}
        verseLabel={`Quran ${data.number}:${verse.number}`}
      />

      <section data-chapter-verse={`${data.number}:${verse.number}`}>
        {mode !== 'english' && <VerseArabic verse={verse} wordByWord={wordByWord} transliteration={showTransliteration} />}
        {mode !== 'arabic' && <p className='text-base sm:text-lg leading-relaxed'>{verse.text}</p>}
        {showCommentary && verse.comments && <div className='mt-3 p-3 bg-yellow-50 text-sm leading-relaxed'><strong>Tafsir (stored commentary):</strong> {verse.comments}</div>}
      </section>
      <Link className='my-4 text-sm underline underline-offset-4' to={buildQuranShareUrl(`/quran/${data.number}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration }) + `#${verse.number}`}>Read full chapter at {data.number}:{verse.number}</Link>
      <div className='flex justify-between mt-6 sm:mt-8 mb-8 sm:mb-12'>
        {hasPreviousVerse ? (
          <PrimaryButton asChild disabled={!hasPreviousVerse} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${chapterNumber}/${previousVerse}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasPreviousChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${previousChapter}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        )}

        {hasNextVerse ? (
          <PrimaryButton asChild disabled={!hasNextVerse} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${chapterNumber}/${nextVerse}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasNextChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${nextChapter}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        )}
      </div>
      <p className='text-xs text-gray-500 mt-4'>Arabic: <a className='underline' href='https://quranenc.com/en/home'>QuranEnc · Uthmani</a></p>
    </div>
  );
}
