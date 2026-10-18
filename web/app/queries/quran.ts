import { useMemo } from 'react';
import { useSuspenseQuery } from '@tanstack/react-query';
import type { Translation } from '~/components/quran/reading-settings';
import { httpGet } from '~/utils/http';

type ChapterType = {
  english: string;
  name: string;
  number: number;
};

type ListSurahResponseType = ChapterType[];

export const listSurahsOptions = () => ({
  queryKey: ['list-surahs'],
  queryFn: async () => {
    const data = await httpGet<ListSurahResponseType>('/api/quran/chapters');
    return data;
  },
});

type WordByWord = {
  arabic: string;
  english: string;
  transliteration: string;
};

type VerseType = {
  chapter: number;
  number: number;
  text: string;
  translations?: Partial<Record<Translation, string>>;
  arabic: string;
  words?: WordByWord[];
  comments?: string;
  audio_arabic?: string;
  audio_english?: string;
};

type ChapterResponseType = {
  name: string;
  english: string;
  number: number;
  verses: VerseType[];
};

export const getChapterOptions = (chapterNumber: number) => ({
  queryKey: ['get-chapter', chapterNumber],
  queryFn: async () => {
    const data = await httpGet<ChapterResponseType>(
      `/api/quran/${chapterNumber}`
    );
    return data;
  },
});

// Both editions share the same chapter cache, so switching also works offline.
export function useTranslatedChapter(chapterNumber: number, translation: Translation) {
  const { data: original } = useSuspenseQuery(getChapterOptions(chapterNumber));
  return useMemo(() => ({
    ...original,
    verses: original.verses.map(verse => ({
      ...verse,
      text: translation === 'khattab' ? verse.text :
        verse.translations?.[translation] ?? 'This translation has not been saved yet. Reconnect to update this chapter.',
    })),
  }), [original, translation]);
}
