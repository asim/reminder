import { httpGet } from '~/utils/http';

type BookType = {
  name: string;
  number: number;
  hadith_count: number;
};

type ListBooksResponseType = BookType[];

export const listBooksOptions = () => ({
  queryKey: ['list-hadith-books'],
  queryFn: async () => {
    const data = await httpGet<ListBooksResponseType>('/api/hadith/books');
    return data;
  },
});

type HadithType = {
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

export const getBookOptions = (bookNumber: number) => ({
  queryKey: ['get-hadith-book', bookNumber],
  queryFn: async () => {
    const data = await httpGet<BookResponseType>(`/api/hadith/${bookNumber}`);
    return data;
  },
});
