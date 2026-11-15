import { readingMeta } from '~/utils/reading-meta';
import { readingData } from '~/utils/prerender.server';
import type { Route } from '.react-router/types/app/routes/+types/_app.quran.$chapterNumber.$verseNumber';
import { CircleChevronLeft, CircleChevronRight } from 'lucide-react';
import { Link } from 'react-router';
import { EmbedButton } from '~/components/interface/embed-button';
import { PageError } from '~/components/interface/page-error';
import { PrimaryButton } from '~/components/interface/primary-button';
import { ShareButton } from '~/components/interface/share-button';
import { AudioPlayer } from '~/components/quran/audio-player';
import { ChapterHeader } from '~/components/quran/chapter-header';
import { VerseArabic } from '~/components/quran/verse-arabic';
import { queryClient } from '~/utils/query-client';
import { getVersePageOptions, useTranslatedChapter } from '~/queries/quran';
import { useQuranSettings } from '~/components/quran/reading-settings';
import { buildQuranShareUrl } from '~/utils/quran-share';

export default function QuranVerse(props: Route.ComponentProps) {
  const { chapterNumber, verseNumber } = props.params;

  const { settings, textStyle } = useQuranSettings();
  const { mode, wordByWord, transliteration: showTransliteration, commentary: showCommentary, translation } = settings;
  const data = useTranslatedChapter(Number(chapterNumber), translation, Number(verseNumber));

  const verse = data.verses.find(
    (verse) => verse.number === Number(verseNumber)
  );

  if (!verse) {
    return <PageError error={new Error('Verse not found')} />;
  }

  const previousVerse = Number(verseNumber) - 1;
  const nextVerse = Number(verseNumber) + 1;
  const hasPreviousVerse = previousVerse >= 1;

  // Use the complete chapter count even when only this verse is hydrated.
  const lastVerse = data.verse_count;
  const hasNextVerse = nextVerse <= lastVerse;

  const nextChapter = Number(chapterNumber) + 1;
  const previousChapter = Number(chapterNumber) - 1;
  const hasNextChapter = nextChapter <= 114;
  const hasPreviousChapter = previousChapter >= 1;

  return (
    <div style={textStyle} className='quran-reader max-w-4xl flex flex-col w-full mx-auto p-0 lg:p-8'>

      <ChapterHeader
        title={data.name}
        translation={data.english}
        subtitle={`Verse ${data.number}:${verseNumber}`}
        actions={<>
        <EmbedButton path={`/quran/${data.number}/${verseNumber}?translation=${translation}`} />
        <ShareButton
          title={`Quran ${data.number}:${verseNumber} - ${data.english}`}
          text={verse.text}
          url={buildQuranShareUrl(`/quran/${data.number}/${verseNumber}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}
        />
        </>}
      />

      {/* Audio Player */}
      <AudioPlayer
        arabicUrl={verse.audio_arabic}
        englishUrl={verse.audio_english}
        verseLabel={`Quran ${data.number}:${verse.number}`}
      />

      <section data-chapter-verse={`${data.number}:${verse.number}`}>
        {mode !== 'english' && <VerseArabic verse={verse} wordByWord={wordByWord} transliteration={showTransliteration} />}
        {mode !== 'arabic' && <p className='quran-translation text-base sm:text-lg leading-relaxed'>{verse.text}</p>}
        {showCommentary && verse.comments && <div className='mt-3 p-3 bg-yellow-50 text-sm leading-relaxed'><strong>{verse.commentary_source || "Stored commentary"}:</strong> {verse.comments}</div>}
      </section>
      <Link className='my-4 text-sm underline underline-offset-4' to={buildQuranShareUrl(`/quran/${data.number}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration }) + `#${verse.number}`}>Read full chapter at {data.number}:{verse.number}</Link>
      <div className='flex justify-between mt-6 sm:mt-8 mb-8 sm:mb-12'>
        {hasPreviousVerse ? (
          <PrimaryButton asChild disabled={!hasPreviousVerse} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${chapterNumber}/${previousVerse}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasPreviousChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${previousChapter}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <CircleChevronLeft className='size-3 sm:size-4' />
              <span className='ml-1'>Previous</span>
            </Link>
          </PrimaryButton>
        )}

        {hasNextVerse ? (
          <PrimaryButton asChild disabled={!hasNextVerse} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${chapterNumber}/${nextVerse}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        ) : (
          <PrimaryButton asChild disabled={!hasNextChapter} className='text-sm sm:text-base py-1 sm:py-2 px-2 sm:px-4'>
            <Link to={buildQuranShareUrl(`/quran/${nextChapter}`, { mode, translation, wordByWord, commentary: showCommentary, transliteration: showTransliteration })}>
              <span className='mr-1'>Next</span>
              <CircleChevronRight className='size-3 sm:size-4' />
            </Link>
          </PrimaryButton>
        )}
      </div>
    </div>
  );
}

export function loader({ params }: { params: Record<string, string | undefined> }) {
  const chapter = readingData('quran', params.chapterNumber);
  const verses = chapter.verses.filter((v: { number: number }) => v.number === Number(params.verseNumber));
  if (!verses.length) throw new Response('Not found', { status: 404 });
  return { readingQueries: [{ queryKey: ['get-verse-page', Number(params.chapterNumber), Number(params.verseNumber)], data: { ...chapter, verses } }] };
}

export const meta = (args: any) => readingMeta('verse', args);

export async function clientLoader({ params }: { params: Record<string, string | undefined> }) {
  const options = getVersePageOptions(Number(params.chapterNumber), Number(params.verseNumber));
  const data = await queryClient.ensureQueryData(options);
  return { readingQueries: [{ queryKey: options.queryKey, data }] };
}
