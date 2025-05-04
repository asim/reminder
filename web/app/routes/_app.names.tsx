import { useSuspenseQuery } from '@tanstack/react-query';
import { Search } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router';
import { listNamesOptions } from '~/queries/names';
import { queryClient } from '~/utils/query-client';

export async function clientLoader() {
  await queryClient.ensureQueryData(listNamesOptions());
}

export default function Names() {
  const { data: names } = useSuspenseQuery(listNamesOptions());
  const [search, setSearch] = useState('');

  const location = useLocation();

  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [search, location.pathname]);

  return (
    <div className='flex flex-row h-full'>
      <div className='flex flex-col pb-15 w-[250px] border-r border-gray-200 overflow-y-auto'>
        <div className='pl-7 pr-4 border-b border-gray-200 relative sticky top-0 bg-white'>
          <Search
            size={15}
            className='absolute left-5 top-1/2 -translate-y-1/2 text-gray-500'
          />
          <input
            type='text'
            placeholder='Search names'
            className='px-4 py-2 focus:outline-none text-sm'
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {names
          ?.filter(
            (name) =>
              !search ||
              name.meaning.toLowerCase().includes(search.toLowerCase()) ||
              name.english.toLowerCase().includes(search.toLowerCase())
          )
          .map((name) => (
            <NavLink
              key={name.number}
              to={`/names/${name.number}`}
              className={({ isActive }) =>
                `py-2 px-4 text-sm border-b border-gray-200 hover:bg-gray-50 cursor-pointer ${
                  isActive ? 'bg-black pointer-events-none text-white' : ''
                }`
              }
            >
              <span className='tabular-nums w-[20px] text-right mr-3 text-xs text-gray-400 inline-block'>
                {name.number}
              </span>
              {name.meaning}
            </NavLink>
          ))}
      </div>

      <div
        ref={containerRef}
        className='flex flex-col overflow-y-auto flex-1 px-5 py-5'
      >
        <Outlet />
      </div>
    </div>
  );
}
