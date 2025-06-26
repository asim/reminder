import { Link } from 'react-router';
import type { Route } from './+types/_index';
import { Code, Globe2 } from 'lucide-react';
import { useEffect, useState } from 'react';

function getHijriDate(): { year: number; month: string; day: number } {
  // Umm al-Qura calendar approximation
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

  // Accurate Hijri conversion (Tabular Islamic Calendar)
  // Source: https://webspace.science.uu.nl/~gent0113/islam/addfiles/islamtabcal.htm
  const jd =
    Math.floor((1461 * (year + 4800 + Math.floor((month - 14) / 12))) / 4) +
    Math.floor((367 * (month - 2 - 12 * Math.floor((month - 14) / 12))) / 12) -
    Math.floor((3 * Math.floor((year + 4900 + Math.floor((month - 14) / 12)) / 100)) / 4) +
    day -
    32075;
  const islamicEpoch = 1948439;
  const days = jd - islamicEpoch;
  let hYear = Math.floor((30 * days + 10646) / 10631);
  let hMonth = Math.min(
    12,
    Math.ceil((days - 29 - hijriToJD(hYear, 1, 1)) / 29.5) + 1
  );
  let hDay = jd - hijriToJD(hYear, hMonth, 1) + 1;

  function hijriToJD(year: number, month: number, day: number): number {
    return (
      day +
      Math.ceil(29.5 * (month - 1)) +
      (year - 1) * 354 +
      Math.floor((3 + 11 * year) / 30) +
      islamicEpoch -
      1
    );
  }

  // Clamp values to valid ranges
  hMonth = Math.max(1, Math.min(12, hMonth));
  hDay = Math.max(1, Math.min(30, hDay));

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
        <div className='font-semibold text-lg'>Salam,</div>
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
