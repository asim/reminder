import { useQuery } from '@tanstack/react-query';
import { useParams } from 'react-router';
import { httpPost } from '~/utils/http';

interface DailyResponse {
  name: string;
  hadith: string;
  verse: string;
  links: Record<string, string>;
  updated: string;
  message: string;
  date: string;
  hijri?: string;
}

function formatDate(dateString: string) {
  if (dateString == "latest") {
    return "Latest Reminder";
  }

  const date = new Date(dateString);
  // Options for formatting the date
  const options: Intl.DateTimeFormatOptions = {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  };
  return date.toLocaleDateString('en-GB', options);
}

export default function DailyByDate() {
  const { date } = useParams<{ date: string }>();
  const { data, isLoading, error } = useQuery<DailyResponse>({
    queryKey: ['daily-by-date', date],
    queryFn: async () => {
      if (!date) throw new Error('No date provided');
      return httpPost<DailyResponse>(`/api/daily/${date}`);
    },
    enabled: !!date,
  });

  if (!date) return <div className="p-4 text-red-500">No date provided</div>;
  if (isLoading) return <div className="p-4">Loading...</div>;
  if (error || !data) return <div className="p-4 text-red-500">Failed to load daily reminder. Please check if the server is running.</div>;

  // Defensive: extract fields, but render whatever is present
  const { verse, hadith, name, links: rawLinks, updated, message, hijri: hijriDate } = data;
  const links = rawLinks || {};

  return (
    <div className="max-w-4xl mx-auto w-full mb-8 sm:mb-12 flex-grow p-0 lg:p-8 space-y-8">
      <div className="text-center">
        <div className="text-base sm:text-lg md:text-xl text-gray-600">{data.hijri}</div>
        <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold mb-1 sm:mb-2">{formatDate(data.date)}</h1>
      </div>

      {message && (
        <section>
          <div className="whitespace-pre-wrap leading-relaxed bg-purple-50 rounded-lg p-4 sm:p-6 text-sm sm:text-base shadow-sm border border-purple-100">
            {message}
          </div>
        </section>
      )}
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Verse</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A verse from the Quran</div>
        <div className="whitespace-pre-wrap leading-snug bg-blue-50 rounded p-4 text-base shadow">
          {links.verse ? <a href={links.verse}>{verse}</a> : verse}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Hadith</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A hadith from sahih bukhari</div>
        <div className="whitespace-pre-wrap leading-snug bg-green-50 rounded p-4 text-base shadow">
          {links.hadith ? <a href={links.hadith}>{hadith}</a> : hadith}
        </div>
      </section>
      <section>
        <h2 className="text-lg font-semibold mb-2 mt-2">Name of Allah</h2>
        <div className="text-sm sm:text-base text-gray-700 mb-2">A beautiful name from the 99 names of Allah</div>
        <div className="whitespace-pre-wrap leading-snug bg-yellow-50 rounded p-4 text-base shadow">
          {links.name ? <a href={links.name}>{name}</a> : name}
        </div>
      </section>
      <section>
        <div className="mt-2">Updated {data.updated}</div>
      </section>
    </div>
  );
}
