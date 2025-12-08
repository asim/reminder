import { useSuspenseQuery } from '@tanstack/react-query';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';
import { httpGet } from '~/utils/http';

interface DailyIndexEntry {
  verse: string;
  hadith: string;
  name: string;
  date: string; // hijri
  gregorian: string;
  message: string;
  hijri?: string;
}

function formatDate(dateString: string) {
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

export default function DailySidebarNav() {
  const { data } = useSuspenseQuery({
    queryKey: ['daily-index'],
    queryFn: async () => httpGet<Record<string, DailyIndexEntry>>('/api/daily/index'),
  });
  if (!data) return null;

  // Filter out 'latest' entry and sort by date descending
  const entries = Object.entries(data)
    .filter(([date]) => date !== 'latest')
    .sort((a, b) => b[0].localeCompare(a[0])) as [string, DailyIndexEntry][];

  const sidebarItems: SidebarItem[] = entries.map(([date, entry], index) => ({
    key: date,
    text: formatDate(entry.date),
    path: `/daily/${date}`,
    number: index + 1,
    extra: (entry as DailyIndexEntry).hijri,
    searchableText: [entry.date, entry.verse, entry.hadith, entry.name, entry.message],
  }));

  return (
    <SearchableSidebar
      items={sidebarItems}
      searchPlaceholder='Search reminders'
    />
  );
}
