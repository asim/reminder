import { Code, MessageCircle } from 'lucide-react';
import { Link, Outlet } from 'react-router';

export default function AppLayout() {
  return (
    <div className='flex flex-col h-screen'>
      <div className='w-full text-sm justify-center py-2 bg-black text-white flex flex-row items-center gap-1'>
        <Link
          className='bg-white text-black px-2 py-0.5 rounded-md text-xs'
          to='/quran/1'
        >
          Quran
        </Link>
        <Link
          className='bg-white text-black px-2 py-0.5 rounded-md text-xs'
          to='/hadith'
        >
          Hadith
        </Link>
        <Link
          className='bg-white text-black px-2 py-0.5 rounded-md text-xs'
          to='/names'
        >
          Names of Allah
        </Link>

        <span className='text-gray-400 mx-2'>/</span>
        <Link
          className='bg-white text-black px-2 py-0.5 rounded-md text-xs flex flex-row items-center gap-1'
          to='/quran/1'
        >
          <MessageCircle className='size-3' />
          Chat with AI
        </Link>
        <Link
          className='bg-white text-black px-2 py-0.5 rounded-md text-xs flex flex-row items-center gap-1'
          to='/quran/1'
        >
          <Code className='size-3' />
          Use API
        </Link>
      </div>
      <Outlet />
    </div>
  );
}
