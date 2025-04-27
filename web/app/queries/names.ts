import { httpGet } from '~/utils/http';

type NameType = {
  number: number;
  english: string;
  arabic: string;
  meaning: string;
  description: string;
  summary: string;
  location: string[];
};

type ListNamesResponseType = NameType[];

export const listNamesOptions = () => ({
  queryKey: ['list-names'],
  queryFn: async () => {
    const data = await httpGet<ListNamesResponseType>('/api/names');
    return data;
  },
});

export const getNameOptions = (nameNumber: number) => ({
  queryKey: ['get-name', nameNumber],
  queryFn: async () => {
    const data = await httpGet<NameType>(`/api/names/${nameNumber}`);
    return data;
  },
});
