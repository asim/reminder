import type { Route } from '.react-router/types/app/routes/+types/_app.hadith';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import { Outlet, useLocation } from 'react-router';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';
import { listBooksOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(listBooksOptions());
}

export default function Hadith() {
  const { data: books } = useSuspenseQuery(listBooksOptions());
  const location = useLocation();
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [location.pathname]);

  const sidebarItems: SidebarItem[] = books?.map((book, counter) => {
    const number = counter + 1;
    return {
      key: book.name,
      text: book.name,
      path: `/hadith/${number}`,
      number,
      searchableText: [book.name, book.english],
      extra: book.english
    };
  }) || [];

  return (
    <div className='flex flex-row h-full'>
      <SearchableSidebar
        items={sidebarItems}
        searchPlaceholder="Search books"
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
