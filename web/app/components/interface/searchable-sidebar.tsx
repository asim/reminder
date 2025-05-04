import { Search } from 'lucide-react';
import { useState } from 'react';
import { NavLink } from 'react-router';

export type SidebarItem = {
  key: string | number;
  text: string;
  path: string;
  number: string | number;
  searchableText: string[];
};

type SearchableSidebarProps = {
  items: SidebarItem[];
  searchPlaceholder: string;
};

export function SearchableSidebar(props: SearchableSidebarProps) {
  const { items, searchPlaceholder } = props;
  const [search, setSearch] = useState('');

  return (
    <div className='flex flex-col pb-15 w-[250px] border-r border-gray-200 overflow-y-auto'>
      <div className='pl-7 pr-4 border-b border-gray-200 relative sticky top-0 bg-white'>
        <Search
          size={15}
          className='absolute left-5 top-1/2 -translate-y-1/2 text-gray-500'
        />
        <input
          type='text'
          placeholder={searchPlaceholder}
          className='px-4 py-2 focus:outline-none text-sm'
          value={search}
          onChange={(e) => setSearch(e.target.value)}
        />
      </div>

      {items
        ?.filter((item) => {
          if (!search) return true;
          const searchLower = search.toLowerCase();
          return item.searchableText.some((text) =>
            text.toLowerCase().includes(searchLower)
          );
        })
        .map((item) => (
          <NavLink
            key={item.key}
            to={item.path}
            className={({ isActive }) =>
              `py-2 px-4 flex items-center text-sm border-b border-gray-200 hover:bg-gray-50 cursor-pointer ${
                isActive ? 'bg-black pointer-events-none text-white' : ''
              }`
            }
          >
            <span className='tabular-nums w-[20px] text-right mr-3 text-xs text-gray-400 inline-block'>
              {item.number}
            </span>
            <span className='truncate'>{item.text}</span>
          </NavLink>
        ))}
    </div>
  );
} 