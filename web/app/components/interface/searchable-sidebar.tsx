import { MenuIcon, Search } from 'lucide-react';
import { useEffect, useState } from 'react';
import { NavLink, useLocation } from 'react-router';
import { cn } from '~/utils/classname';

export type SidebarItem = {
  key: string | number;
  text: string;
  path: string;
  number: string | number;
  searchableText: string[];
  extra?: string;
};

type SearchableSidebarProps = {
  items: SidebarItem[];
  searchPlaceholder: string;
};

export function SearchableSidebar(props: SearchableSidebarProps) {
  const { items, searchPlaceholder } = props;
  const [search, setSearch] = useState('');

  const [isOpen, setIsOpen] = useState(false);
  const location = useLocation();

  useEffect(() => {
    setIsOpen(false);
  }, [location.pathname]);

  return (
    <>
      <div className='px-2 block lg:hidden border-r border-gray-200 py-2 text-sm'>
        <button
          className='flex items-center gap-2 text-gray-500 hover:text-black'
          onClick={() => setIsOpen(!isOpen)}
        >
          <MenuIcon className='size-4' />
        </button>
      </div>
      <div
        onClick={() => setIsOpen(false)}
        className={cn('fixed inset-0 z-40 bg-black/50 ', {
          'opacity-100 pointer-events-auto': isOpen,
          'opacity-0 pointer-events-none': !isOpen,
        })}
      />
      <div
        className={cn(
          'hidden lg:flex flex-col pb-15 w-[250px] border-r border-gray-200 overflow-y-auto',
          {
            'fixed flex bg-white z-50 top-0 left-0 bottom-0': isOpen,
          }
        )}
      >
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
              title={item.text}
              className={({ isActive }) =>
                `py-2 px-4 flex items-center text-sm border-b border-gray-200 hover:bg-gray-50 cursor-pointer ${
                  isActive ? 'bg-black pointer-events-none text-white' : ''
                }`
              }
            >
              <span className='tabular-nums w-[28px] shrink-0 text-right mr-3 text-xs text-gray-400 inline-block'>
                {item.number}
              </span>
              <div className="truncate min-w-0 flex-1">
                <div className="truncate">{item.text}</div>
                {item.extra && (
                  <div className="text-gray-400 text-xs truncate">{item.extra}</div>
                )}
              </div>
            </NavLink>
          ))}
      </div>
    </>
  );
}
