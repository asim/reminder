import { useQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import React from 'react';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';

interface DailyIndexEntry {
  verse: string;
  hadith: string;
  name: string;
  date: string; // hijri
  gregorian: string;
  message: string;
}

export default function DailySidebar() {
  const { data, isLoading, error } = useQuery<Record<string, DailyIndexEntry>>({
    queryKey: ['daily-index'],
    queryFn: async () => httpGet<Record<string, DailyIndexEntry>>('/api/daily/index'),
  });

  if (isLoading) return <div className='p-2 text-xs'>Loading daily index...</div>;
  if (error) return <div className='p-2 text-xs text-red-500'>Failed to load daily index.</div>;
  if (!data) return null;

  // Sort by date descending
  const entries = Object.entries(data).sort((a, b) => b[0].localeCompare(a[0]));

  const sidebarItems: SidebarItem[] = entries.map(([date, entry], index) => ({
    key: date,
    text: entry.hijri,
    path: `/daily/${date}`,
    number: index + 1, // Use index as a simple number
    extra: entry.date,
    searchableText: [entry.hijri, entry.date, entry.verse, entry.hadith, entry.name, entry.message],
  }));

  return (
    <SearchableSidebar
      items={sidebarItems}
      searchPlaceholder='Search daily reminders'
    />
  );
}
