import { Link } from 'react-router';

export default function App() {
  return (
    <div className='flex flex-col items-center justify-center min-h-screen px-4 py-16 bg-white'>
      <h1 className='text-7xl mb-4 font-bold tracking-tight'>reminder</h1>
      <p className='text-gray-600 text-xl mb-10 text-center'>
        Quran, hadith, and more as an app and API
      </p>

      <div className='flex flex-col gap-2 w-full max-w-md'>
        <Link
          to='/app/quran'
          className='bg-white flex flex-col gap-3 items-center justify-center text-black p-3 w-full rounded-lg border border-gray-200 hover:border-gray-400 transition-all duration-200 cursor-pointer'
        >
          Quran
        </Link>
        <Link
          to='/app/hadith'
          className='bg-white flex flex-col gap-3 items-center justify-center text-black p-3 w-full rounded-lg border border-gray-200 hover:border-gray-400 transition-all duration-200 cursor-pointer'
        >
          Hadith
        </Link>
        <Link
          to='/app/names'
          className='bg-white flex flex-col gap-3 items-center justify-center text-black p-3 w-full rounded-lg border border-gray-200 hover:border-gray-400 transition-all duration-200 cursor-pointer'
        >
          Allah Names
        </Link>
        <Link
          to='/app/search'
          className='bg-white flex flex-col gap-3 items-center justify-center text-black p-3 w-full rounded-lg border border-gray-200 hover:border-gray-400 transition-all duration-200 cursor-pointer'
        >
          Search
        </Link>
      </div>
    </div>
  );
}
