import { Code, Globe2, Search } from 'lucide-react';
import { Link, NavLink } from 'react-router';
import type { Route } from './+types/_index';

export function meta({ }: Route.MetaArgs) {
  return [
    { title: 'Reminder - Quran & Hadith' },
    {
      name: 'description',
      content: 'Access Quran, hadith, and more through our app and API',
    },
  ];
}

export default function Home() {
  return (
    <div className='flex flex-col items-center justify-center min-h-screen px-4 sm:px-6 md:px-8 py-12 sm:py-16 bg-white'>
      {/* Fixed nav bar at the top, matching AppLayout */}
      <nav className='fixed top-0 left-0 w-full z-20 bg-black text-white flex flex-row items-center gap-1 px-2 py-2 text-sm justify-center'>
        <NavLink to='/daily' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1'>Daily</NavLink>
        <NavLink to='/quran' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1'>Quran</NavLink>
        <NavLink to='/hadith' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1'>Hadith</NavLink>
        <NavLink to='/names' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1'>Names</NavLink>
        <NavLink to='/islam' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1'>Islam</NavLink>
        <NavLink to='/search' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1 flex items-center gap-1'><Search className='size-3' />Search</NavLink>
        <NavLink to='/api' className='bg-white text-black px-2 py-0.5 rounded-md text-xs mx-1 flex items-center gap-1'><Code className='size-3' />API</NavLink>
      </nav>
      <div className='pt-16 w-full flex flex-col items-center'>
        <h1 className='text-4xl sm:text-5xl md:text-6xl lg:text-7xl mb-2 sm:mb-4 font-bold tracking-tight text-center'>reminder</h1>
        <p className='text-gray-600 text-balance text-lg sm:text-xl mb-6 sm:mb-8 md:mb-10 text-center max-w-md'>
          Quran, hadith, and more as an app and API
        </p>

        <div className='flex flex-row lg:flex-col gap-2 w-full max-w-xs sm:max-w-sm md:max-w-md'>
          <Link
            to='/quran'
            className='flex flex-row gap-2 sm:gap-3 items-center justify-start py-2.5 sm:py-3 px-4 sm:px-5 w-full rounded-lg border border-black bg-black text-white hover:opacity-70 transition-all duration-200 cursor-pointer'
          >
            <Globe2 className='size-4' />
            <span className='font-medium mr-1'>App</span>
            <span className='hidden lg:inline'>Read Quran, hadith, and more</span>
          </Link>
          <Link
            to='/api'
            className='bg-white flex flex-row gap-2 sm:gap-3 items-center justify-start text-black py-2.5 sm:py-3 px-4 sm:px-5 w-full rounded-lg border border-gray-500 hover:bg-gray-200 hover:border-gray-600 hover:text-black hover:border-black transition-all duration-200 cursor-pointer'
          >
            <Code className='size-4' />
            <span className='font-medium mr-1'>API</span>
            <span className='hidden lg:inline'>Develop using our free API</span>
          </Link>
        </div>
      </div>
    </div>
  );
}
