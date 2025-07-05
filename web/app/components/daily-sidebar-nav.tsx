import { useQuery } from '@tanstack/react-query';
import { httpGet } from '~/utils/http';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';
import React from 'react';

interface DailyIndexEntry {
  verse: string;
  hadith: string;
  name: string;
  date: string; // hijri
  gregorian: string;
  message: string;
}

function formatDate(dateString) {
  const date = new Date(dateString); // Create a Date object from the "YYYY-MM-DD" string

  // Options for formatting the date
  const options = {
    weekday: 'long', // e.g., "Saturday"
    year: 'numeric', // e.g., "2025"
    month: 'long',   // e.g., "June"
    day: 'numeric'   // e.g., "25"
  };

  // Use toLocaleDateString to format the date
  // The 'en-GB' locale provides day before month, and the options tailor the output
  return date.toLocaleDateString('en-GB', options);
}

export default function DailySidebarNav() {
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
    text: formatDate(entry.date),
    path: `/daily/${date}`,
    number: index + 1,
    extra: entry.hijri,
    searchableText: [entry.date, entry.verse, entry.hadith, entry.name, entry.message],
  }));

  return (
    <SearchableSidebar
      items={sidebarItems}
      searchPlaceholder='Search daily reminders'
    />
  );
}
