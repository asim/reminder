import type { Route } from '.react-router/types/app/routes/+types/_app.quran';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import { Outlet, useLocation } from 'react-router';
import {
  SearchableSidebar,
  type SidebarItem,
} from '~/components/interface/searchable-sidebar';
import { listSurahsOptions } from '~/queries/quran';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(listSurahsOptions());
}

export default function Quran() {
  const { data: chapters } = useSuspenseQuery(listSurahsOptions());
  const location = useLocation();
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [location.pathname]);

  const sidebarItems: SidebarItem[] =
    chapters?.map((chapter) => ({
      key: chapter.number,
      text: chapter.name,
      path: `/quran/${chapter.number}`,
      number: chapter.number,
      searchableText: [chapter.name],
    })) || [];

  return (
    <div className='flex flex-row h-full'>
      <SearchableSidebar
        items={sidebarItems}
        searchPlaceholder='Search surah'
      />

      <div
        ref={containerRef}
        className='flex flex-col overflow-y-auto flex-1 px-5 py-5'
      >
        <Outlet />
      </div>
    </div>
  );
}
