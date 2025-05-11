import { useSuspenseQuery } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import { Outlet, useLocation } from 'react-router';
import { SearchableSidebar, type SidebarItem } from '~/components/interface/searchable-sidebar';
import { listNamesOptions } from '~/queries/names';
import { queryClient } from '~/utils/query-client';

export async function clientLoader() {
  await queryClient.ensureQueryData(listNamesOptions());
}

export default function Names() {
  const { data: names } = useSuspenseQuery(listNamesOptions());
  const location = useLocation();
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [location.pathname]);

  const sidebarItems: SidebarItem[] = names?.map(name => ({
    key: name.number,
    text: name.english,
    path: `/names/${name.number}`,
    number: name.number,
    searchableText: [name.meaning, name.english],
    extra: name.meaning
  })) || [];

  return (
    <div className='flex flex-row h-full'>
      <SearchableSidebar
        items={sidebarItems}
        searchPlaceholder="Search names"
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
