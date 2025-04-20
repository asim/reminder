import { Link } from 'react-router';
import type { Route } from './+types/_index';
import { Code, Globe2 } from 'lucide-react';

export function meta({}: Route.MetaArgs) {
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
    <div className='flex flex-col items-center justify-center min-h-screen px-4 py-16 bg-white'>
      <h1 className='text-7xl mb-4 font-bold tracking-tight'>reminder</h1>
      <p className='text-gray-600 text-xl mb-10 text-center'>
        Quran, hadith, and more as an app and API
      </p>

      <div className='flex flex-col gap-2 w-full max-w-md'>
        <Link
          to='/app'
          className='flex flex-row gap-3 items-center justify-start py-3 px-5 w-full rounded-lg border border-black bg-black text-white hover:opacity-70 transition-all duration-200 cursor-pointer'
        >
          <Globe2 className='size-4' />
          <span className='font-medium mr-1'>Application</span> Read Quran,
          hadith, and more
        </Link>
        <Link
          to='/api'
          className='bg-white flex flex-row gap-3 items-center justify-start text-black py-3 px-5 w-full rounded-lg border border-gray-500 hover:bg-gray-200 hover:border-gray-600 hover:text-black hover:border-black transition-all duration-200 cursor-pointer'
        >
          <Code className='size-4' />
          <span className='font-medium mr-1'>API</span> Develop using our free
          API
        </Link>
      </div>
    </div>
  );
}
