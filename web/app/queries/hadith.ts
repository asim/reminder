import { hadithPath, type HadithCollection } from '~/utils/hadith-collection';
import { httpGet } from '~/utils/http';

type BookType = {
  name: string;
  number: number;
  hadith_count: number;
};

type ListBooksResponseType = BookType[];

export const listBooksOptions = (collection: HadithCollection = 'bukhari') => ({
  queryKey: ['list-hadith-books', collection],
  queryFn: async () => {
    const data = await httpGet<ListBooksResponseType>(`/api${hadithPath(collection)}/books`);
    return data;
  },
});

type HadithType = {
  number: number;
  narrator: string;
  english: string;
  arabic: string;
  chain?: string;
  // Legacy fields
  by: string;
  info: string;
  text: string;
};

type BookResponseType = {
  name: string;
  number: number;
  hadiths: HadithType[];
  hadith_count: number;
};

export const getBookOptions = (bookNumber: number, collection: HadithCollection = 'bukhari') => ({
  queryKey: ['get-hadith-book', collection, bookNumber],
  queryFn: async () => {
    const data = await httpGet<BookResponseType>(`/api${hadithPath(collection)}/${bookNumber}`);
    return data;
  },
});
