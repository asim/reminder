import { Code, Search } from 'lucide-react';
import { Outlet, NavLink } from 'react-router';
import { cn } from '~/utils/classname';

export default function AppLayout() {
  const buttonClass = ({ isActive }: { isActive: boolean }) =>
    cn(
      'bg-white transition-colors hover:opacity-100 flex items-center gap-1 text-black px-2 py-0.5 rounded-md text-xs',
      isActive
        ? 'bg-white opacity-100 text-black'
        : 'bg-black text-white hover:bg-white hover:text-black'
    );

  return (
    <div className='flex flex-col h-screen overflow-hidden'>
      <div className='w-full text-sm justify-center py-2 px-2 bg-black text-white flex flex-row items-center gap-1'>
        <NavLink className={buttonClass} to='/quran'>
          Quran
        </NavLink>
        <NavLink className={buttonClass} to='/hadith'>
          Hadith
        </NavLink>
        <NavLink className={buttonClass} to='/names'>
          Names
        </NavLink>
        <NavLink className={buttonClass} to='/daily'>
          Daily
        </NavLink>

        <div className='hidden lg:flex items-center gap-2'>
          <span className='text-gray-400 hidden lg:block mx-1'>/</span>
          <NavLink className={buttonClass} to='/search'>
            <Search className='size-3' />
            Search
          </NavLink>
          <NavLink className={buttonClass} to='/api'>
            <Code className='size-3' />
            API Usage
          </NavLink>
        </div>

        <div className='lg:hidden flex items-center gap-1 ml-auto'>
          <NavLink to='/search' className={buttonClass}>
            <Search className='size-3' />
            Ask
          </NavLink>
          <NavLink to='/api' className={buttonClass}>
            <Code className='size-3' />
            API
          </NavLink>
        </div>
      </div>
      <Outlet />
    </div>
  );
}
