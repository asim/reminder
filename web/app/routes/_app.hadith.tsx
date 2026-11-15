import { HadithSettings } from '~/components/interface/hadith-settings';
import { collectionFromPath, hadithPath } from '~/utils/hadith-collection';
import { readingLoader } from '~/utils/prerender.server';
import type { Route } from '.react-router/types/app/routes/+types/_app.hadith';
import { useSuspenseQuery } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import { Outlet, useLocation } from 'react-router';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';
import { listBooksOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(listBooksOptions(collectionFromPath(new URL(props.request.url).pathname)));
}

export default function Hadith() {
  const location = useLocation();
  const collection = collectionFromPath(location.pathname);
  const { data: books } = useSuspenseQuery(listBooksOptions(collection));
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [location.pathname]);

  const sidebarItems: SidebarItem[] = books?.map((book, counter) => {
    const number = book.number;
    return {
      key: book.name,
      text: book.name,
      path: `${hadithPath(collection)}/${number}`,
      number,
      searchableText: [book.name]
    };
  }) || [];

  return (
    <div className='flex flex-row h-full'>
      <SearchableSidebar
        items={sidebarItems}
        searchPlaceholder="Search books"
        actions={<HadithSettings />}
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

export function loader({ request }: { request: Request }) {
  const collection = collectionFromPath(new URL(request.url).pathname);
  return readingLoader(collection === 'muslim' ? 'muslim-books' : 'books', ['list-hadith-books', collection]);
}
