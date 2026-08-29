import { BookOpen, Code, Download, Moon, Search, Star } from 'lucide-react';
import { Link } from 'react-router';
import type { Route } from './+types/_index';
import { usePwaInstall } from '~/hooks/use-pwa-install';

export function meta({}: Route.MetaArgs) {
  return [
    { title: 'Reminder - Quran, Hadith & Names of Allah' },
    {
      name: 'description',
      content:
        'Read the Quran, Hadith, and learn the Names of Allah. A beautiful, free Islamic app.',
    },
  ];
}

export default function LandingPage() {
  const { canInstall, isInstalled, install } = usePwaInstall();

  return (
    <div className='flex flex-col min-h-screen bg-white'>
      {/* Minimal top bar */}
      <div className='flex items-center justify-between px-4 sm:px-6 py-3'>
        <span className='font-bold text-lg'>Reminder</span>
        <Link
          to='/api'
          className='text-xs sm:text-sm text-gray-500 hover:text-black transition-colors flex items-center gap-1'
        >
          <Code className='size-3.5' />
          Developers
        </Link>
      </div>

      {/* Hero */}
      <div className='flex-1 flex flex-col items-center justify-center px-5 sm:px-6 pb-8 pt-4'>
        <div className='w-full max-w-lg flex flex-col items-center text-center'>
          <div className='mb-6 sm:mb-8 w-16 h-16 sm:w-20 sm:h-20 rounded-2xl bg-black flex items-center justify-center'>
            <span className='text-white text-2xl sm:text-3xl font-bold'>R</span>
          </div>

          <h1 className='text-3xl sm:text-4xl md:text-5xl font-bold tracking-tight mb-3 sm:mb-4'>
            Reminder
          </h1>
          <p className='text-gray-600 text-base sm:text-lg mb-8 sm:mb-10 max-w-sm text-balance leading-relaxed'>
            Read the Quran, Hadith, and learn the beautiful Names of Allah — all in one place.
          </p>

          {/* Install / Open button */}
          <div className='w-full max-w-xs space-y-3 mb-10 sm:mb-12'>
            {canInstall && (
              <button
                onClick={install}
                className='flex items-center justify-center gap-2 w-full py-3 px-5 rounded-xl bg-black text-white font-medium text-base hover:bg-gray-800 transition-colors'
              >
                <Download className='size-4' />
                Install App
              </button>
            )}
            <Link
              to='/home'
              className='flex items-center justify-center gap-2 w-full py-3 px-5 rounded-xl bg-black text-white font-medium text-base hover:bg-gray-800 transition-colors'
            >
              {isInstalled ? 'Open App' : 'Get Started'}
            </Link>
          </div>

          {/* Feature highlights */}
          <div className='grid grid-cols-2 gap-3 sm:gap-4 w-full max-w-sm'>
            <Link
              to='/quran'
              className='flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-gray-400 transition-colors'
            >
              <BookOpen className='size-5 text-gray-700' />
              <span className='text-sm font-medium'>Quran</span>
              <span className='text-xs text-gray-500'>114 Chapters</span>
            </Link>
            <Link
              to='/hadith'
              className='flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-gray-400 transition-colors'
            >
              <Moon className='size-5 text-gray-700' />
              <span className='text-sm font-medium'>Hadith</span>
              <span className='text-xs text-gray-500'>Sahih Bukhari</span>
            </Link>
            <Link
              to='/names'
              className='flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-gray-400 transition-colors'
            >
              <Star className='size-5 text-gray-700' />
              <span className='text-sm font-medium'>Names of Allah</span>
              <span className='text-xs text-gray-500'>99 Names</span>
            </Link>
            <Link
              to='/search'
              className='flex flex-col items-center gap-2 p-4 rounded-xl border border-gray-200 hover:border-gray-400 transition-colors'
            >
              <Search className='size-5 text-gray-700' />
              <span className='text-sm font-medium'>Search</span>
              <span className='text-xs text-gray-500'>AI-powered</span>
            </Link>
          </div>
        </div>
      </div>

      {/* Footer */}
      <div className='py-4 text-center text-xs text-gray-400'>
        Free and open source
      </div>
    </div>
  );
}
