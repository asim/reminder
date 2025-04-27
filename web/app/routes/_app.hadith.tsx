import type { Route } from '.react-router/types/app/routes/+types/_app.hadith';
import { useSuspenseQuery } from '@tanstack/react-query';
import { Search } from 'lucide-react';
import { useEffect, useRef, useState } from 'react';
import { NavLink, Outlet, useLocation } from 'react-router';
import { listBooksOptions } from '~/queries/hadith';
import { queryClient } from '~/utils/query-client';

export async function clientLoader(props: Route.LoaderArgs) {
  await queryClient.ensureQueryData(listBooksOptions());
}

export default function Hadith() {
  const { data: books } = useSuspenseQuery(listBooksOptions());
  const [search, setSearch] = useState('');

  const location = useLocation();

  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (containerRef.current) {
      containerRef.current.scrollTo({ top: 0 });
    }
  }, [search, location.pathname]);

  const booksWithCorrectNumbers = books?.map((book, counter) => ({
    ...book,
    number: counter + 1,
  }));

  return (
    <div className='flex flex-row h-full'>
      <div className='flex flex-col w-[250px] border-r border-gray-200 overflow-y-auto'>
        <div className='pl-7 pr-4 border-b border-gray-200 relative sticky top-0 bg-white'>
          <Search
            size={15}
            className='absolute left-5 top-1/2 -translate-y-1/2 text-gray-500'
          />
          <input
            type='text'
            placeholder='Search books'
            className='px-4 py-2 focus:outline-none text-sm'
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
        </div>

        {booksWithCorrectNumbers
          ?.filter(
            (book) =>
              !search || book.name.toLowerCase().includes(search.toLowerCase())
          )
          .map((book) => (
            <NavLink
              key={book.name}
              to={`/hadith/${book.number}`}
              className={({ isActive }) =>
                `py-2 px-4 text-sm border-b border-gray-200 hover:bg-gray-50 cursor-pointer ${
                  isActive ? 'bg-black pointer-events-none text-white' : ''
                }`
              }
            >
              <span className='tabular-nums w-[20px] text-right mr-3 text-xs text-gray-400 inline-block'>
                {book.number}
              </span>
              {book.name}
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
