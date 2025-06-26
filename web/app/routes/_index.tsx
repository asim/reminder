import { Link } from 'react-router';
import type { Route } from './+types/_index';
import { Code, Globe2 } from 'lucide-react';
import { useEffect, useState } from 'react';

// Pure JS Hijri date conversion (Tabular Islamic calendar)
function getHijriDate() {
  const hijriMonths = [
    'Muharram',
    'Safar',
    'Rabiʿ al-awwal',
    'Rabiʿ al-thani',
    'Jumada al-awwal',
    'Jumada al-thani',
    'Rajab',
    'Shaʿban',
    'Ramadan',
    'Shawwal',
    'Dhu al-Qiʿdah',
    'Dhu al-Hijjah',
  ];
  const today = new Date();
  let day = today.getDate();
  let month = today.getMonth() + 1; // JS months are 0-based
  let year = today.getFullYear();

  // Julian day calculation
  if (month <= 2) {
    year -= 1;
    month += 12;
  }
  const a = Math.floor(year / 100);
  const b = 2 - a + Math.floor(a / 4);
  const jd =
    Math.floor(365.25 * (year + 4716)) +
    Math.floor(30.6001 * (month + 1)) +
    day +
    b -
    1524.5;

  // Hijri date calculation (Tabular Islamic calendar)
  const islamicEpoch = 1948439.5;
  const days = Math.floor(jd - islamicEpoch);
  const hYear = Math.floor((30 * days + 10646) / 10631);
  const firstDayOfYear = islamicEpoch + 354 * (hYear - 1) + Math.floor((3 + 11 * hYear) / 30);
  let hMonth = Math.floor((jd - firstDayOfYear) / 29.5) + 1;
  if (hMonth > 12) hMonth = 12;
  const firstDayOfMonth = firstDayOfYear + 29.5 * (hMonth - 1);
  let hDay = Math.floor(jd - firstDayOfMonth + 1);
  if (hDay < 1) {
    hMonth += 1;
    hDay = 1;
    if (hMonth > 12) {
      hMonth = 1;
    }
  }
  return { year: hYear, month: hijriMonths[hMonth - 1], day: hDay };
}

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
  const [hijri, setHijri] = useState({ year: 0, month: '', day: 0 });

  useEffect(() => {
    setHijri(getHijriDate());
  }, []);

  return (
    <div className='flex flex-col items-center justify-center min-h-screen px-4 sm:px-6 md:px-8 py-12 sm:py-16 bg-white'>
      <div className='absolute top-4 left-4 text-left'>
        <div className='font-semibold text-lg'>Salam</div>
        <div>
          Today is the {hijri.day} of {hijri.month}, {hijri.year}
        </div>
      </div>

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
  );
}
