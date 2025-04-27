import { Code, MessageCircle } from 'lucide-react';
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
    <div className='flex flex-col h-screen'>
      <div className='w-full text-sm justify-center py-2 bg-black text-white flex flex-row items-center gap-1'>
        <NavLink className={buttonClass} to='/quran/1'>
          Quran
        </NavLink>
        <NavLink className={buttonClass} to='/hadith/1'>
          Hadith
        </NavLink>
        <NavLink className={buttonClass} to='/names'>
          Names of Allah
        </NavLink>

        <span className='text-gray-400 mx-2'>/</span>
        <NavLink className={buttonClass} to='/chat'>
          <MessageCircle className='size-3' />
          Chat with AI
        </NavLink>
        <NavLink className={buttonClass} to='/api'>
          <Code className='size-3' />
          Use API
        </NavLink>
      </div>
      <Outlet />
    </div>
  );
}
