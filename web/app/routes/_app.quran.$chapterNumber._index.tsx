import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber._index';
import { useSuspenseQuery } from '@tanstack/react-query';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import React, { useEffect } from 'react';
import { Link, useSearchParams, useLocation } from 'react-router';
import { BookmarkButton } from '~/components/interface/bookmark-button';
import { ReadingBookmarkButton } from '~/components/interface/reading-bookmark-button';
import { ShareButton } from '~/components/interface/share-button';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { AudioPlayer } from '~/components/quran/audio-player';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { VerseArabic } from '~/components/quran/verse-arabic';
import { ReadingControls } from '~/components/quran/reading-controls';
import { ExplainVerse } from '~/components/quran/explain-verse';
import { TranscriptButton } from '~/components/quran/transcript-button';
import { useQuranViewMode } from '~/hooks/use-quran-view-mode';
import { getChapterOptions } from '~/queries/quran';
import { useTransliterationToggle } from '~/use-transliteration-toggle';
import { useWordByWordToggle } from '~/use-word-by-word-toggle';
import { queryClient } from '~/utils/query-client';
import { buildQuranShareUrl } from '~/utils/quran-share';


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
  const location = useLocation();
  const [searchParams] = useSearchParams();
  const [mode, setMode] = useQuranViewMode();
  const [wordByWord, setWordByWord] = useWordByWordToggle();
  const [showTransliteration, setShowTransliteration] = useTransliterationToggle();
  const [showCommentary, setShowCommentary] = React.useState(() => {
    const param = searchParams.get('commentary');
    return param === '1' || param === 'true';
  });
  const [currentVerseIndex, setCurrentVerseIndex] = React.useState(0);
  const [isAutoPlaying, setIsAutoPlaying] = React.useState(false);

  // Filter out verse 0 (Bismillah) for playback
  const playableVerses = data.verses.filter(v => v.number !== 0);

  // Reset playback state when chapter changes
  useEffect(() => {
    setCurrentVerseIndex(0);
    setIsAutoPlaying(false);
  }, [chapterNumber]);

  useEffect(() => {
    if (!data || !window.location.hash) {
      return;
    }

    const timer = window.setTimeout(() => {
      document.getElementById(window.location.hash.substring(1))?.scrollIntoView({ block: 'start' });
    }, 0);
    return () => window.clearTimeout(timer);
  }, [data, location.hash]);

  if (!data) {
    return null;
  }

  const previousChapter = Number(chapterNumber) - 1;
  const nextChapter = Number(chapterNumber) + 1;

  const handleVersePlayComplete = () => {
    // Move to next verse when current verse audio completes
    if (currentVerseIndex < playableVerses.length - 1) {
      setIsAutoPlaying(true);
      setCurrentVerseIndex(currentVerseIndex + 1);
    } else {
      // Reset to first verse when chapter completes
      setIsAutoPlaying(false);
      setCurrentVerseIndex(0);
    }
  };

  const handlePlayStart = () => {
    setIsAutoPlaying(true);
  };

  const currentVerse = playableVerses[currentVerseIndex];

  return (
    <div className='max-w-4xl flex flex-col w-full mb-8 sm:mb-12 flex-grow mx-auto p-0 lg:p-8'>
      <ReadingControls mode={mode} setMode={setMode} wordByWord={wordByWord} setWordByWord={setWordByWord} transliteration={showTransliteration} setTransliteration={setShowTransliteration} commentary={showCommentary} setCommentary={setShowCommentary} />
      <div className='flex justify-end gap-1'>
        <TranscriptButton
          chapterNumber={data.number}
          chapterName={data.name}
          chapterEnglish={data.english}
          verses={data.verses}
        />
        <ShareButton
          title={`${data.english} (${data.name}) - Quran Chapter ${data.number}`}
          url={buildQuranShareUrl(`/quran/${data.number}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}
        />
      </div>
      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Chapter ${data.number}`}
      />
      
      {/* Chapter Audio Player */}
      {currentVerse && (
        <AudioPlayer
          arabicUrl={currentVerse.audio_arabic}
          englishUrl={currentVerse.audio_english}
          verseLabel={`Quran ${data.number}:${currentVerse.number}`}
          onPlayComplete={handleVersePlayComplete}
          onPlayStart={handlePlayStart}
          autoPlay={isAutoPlaying}
        />
      )}
      


      <div className='space-y-6'>
        {data.verses.map(verse => (
          <section id={String(verse.number)} data-chapter-verse={`${data.number}:${verse.number}`} key={verse.number} className='border-b border-gray-100 pb-4 scroll-mt-36'>
            {mode !== 'english' && <VerseArabic verse={verse} wordByWord={wordByWord} transliteration={showTransliteration} />}
            {mode !== 'arabic' && <p className='text-base sm:text-lg leading-relaxed'>{verse.text}</p>}
            {showCommentary && verse.comments && <div className='mt-3 p-3 bg-yellow-50 text-sm leading-relaxed'><strong>Tafsir (stored commentary):</strong> {verse.comments}</div>}
            {verse.number > 0 && <>
              <div className='flex items-center justify-between mt-2 gap-2'>
                <Link className='text-sm underline underline-offset-4' to={buildQuranShareUrl(`/quran/${data.number}/${verse.number}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>{data.number}:{verse.number}</Link>
                <div className='flex items-center gap-1'>
                  <ReadingBookmarkButton type='quran' label={`${data.english} ${data.number}:${verse.number}`} url={`/quran/${data.number}#${verse.number}`} excerpt={verse.text.slice(0, 80)} />
                  <BookmarkButton type='quran' itemKey={`${data.number}:${verse.number}`} label={`${data.english} ${data.number}:${verse.number}`} url={`/quran/${data.number}#${verse.number}`} excerpt={verse.text.slice(0, 80)} />
                  <ShareButton title={`Quran ${data.number}:${verse.number}`} text={verse.text} url={buildQuranShareUrl(`/quran/${data.number}/${verse.number}`, { mode, wordByWord, commentary: showCommentary, transliteration: showTransliteration })} />
                </div>
              </div>
              <ExplainVerse chapter={data.number} verse={verse.number} commentary={verse.comments} />
            </>}
          </section>
        ))}
      </div>
      <div className='flex justify-between mt-6 sm:mt-8 mb-3'>
        <PrimaryButton
          asChild
          disabled={previousChapter < 1}
          className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'
        >
          <Link to={`/quran/${previousChapter}`}>
            <CircleChevronLeft className='size-4' />
            <span className='ml-1'>Previous</span>
          </Link>
        </PrimaryButton>

        <PrimaryButton
          asChild
          disabled={nextChapter > 114}
          className='text-sm sm:text-base py-2 sm:py-2 px-3 sm:px-4'
        >
          <Link to={`/quran/${nextChapter}`}>
            <span className='mr-1'>Next</span>
            <CircleChevronRight className='size-4' />
          </Link>
        </PrimaryButton>
      </div>
      <p className='text-xs text-gray-500 mt-4'>Arabic: <a className='underline' href='https://quranenc.com/en/home'>QuranEnc · Uthmani</a></p>
    </div>
  );
}

export function ErrorBoundary({ error }: Route.ErrorBoundaryProps) {
  return <PageError error={error} />;
}
