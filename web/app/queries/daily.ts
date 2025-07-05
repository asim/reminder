import { httpPost } from '~/utils/http';

export type DailyData = {
  name?: string;
  hadith?: string;
  verse?: string;
  links?: Record<string, string>;
  updated?: string;
  message?: string;
  date?: string;
  hijri?: string;
  error?: string;
};

export const getDailyByDateOptions = (date: string) => ({
  queryKey: ['get-daily', date],
  queryFn: async () => {
    const data = await httpPost<DailyData>(`/api/daily/${date}`);
    return data;
  },
});
