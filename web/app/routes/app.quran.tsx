import type { Route } from '.react-router/types/app/routes/+types/quran';
import { useSuspenseQuery } from '@tanstack/react-query';
import { Search } from 'lucide-react';
import { useState } from 'react';
import { NavLink, Outlet } from 'react-router';
import { listSurahsOptions } from '~/queries/quran';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(listSurahsOptions());
}

export default function Quran() {
  const { data: chapters } = useSuspenseQuery(listSurahsOptions());
  const [search, setSearch] = useState('');

  return (
    <div className='flex flex-row h-screen'>
      <div className='flex flex-col w-[250px] border-r border-gray-200 overflow-y-auto'>
        <div className='pl-7 pr-4 border-b border-gray-200 relative sticky top-0 bg-white'>
          <Search
            size={15}
            className='absolute left-5 top-1/2 -translate-y-1/2 text-gray-500'
          />
          <input
            type='text'
            placeholder='Search surah'
            className='px-4 py-2 focus:outline-none text-sm'
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {chapters
          ?.filter(
            (chapter) =>
              !search ||
              chapter.name.toLowerCase().includes(search.toLowerCase())
          )
          .map((chapter) => (
            <NavLink
              key={chapter.number}
              to={`/app/quran/${chapter.number}`}
              className={({ isActive }) =>
                `py-2 px-4 text-sm border-b border-gray-200 hover:bg-gray-50 cursor-pointer ${
                  isActive ? 'bg-black pointer-events-none text-white' : ''
                }`
              }
            >
              <span className='tabular-nums w-[20px] text-right mr-3 text-xs text-gray-400 inline-block'>
                {chapter.number}
              </span>
              {chapter.name}
            </NavLink>
          ))}
      </div>

      <div className='flex flex-col overflow-y-auto flex-1 px-5 py-5'>
        <Outlet />
      </div>
    </div>
  );
}
