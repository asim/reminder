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

type VerseType = {
  chapter: number;
  number: number;
  text: string;
  arabic: string;
};

type ChapterResponseType = {
  name: string;
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
