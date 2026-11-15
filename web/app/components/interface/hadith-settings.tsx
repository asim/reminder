import { Settings2 } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router';
import { collectionFromPath, hadithPath, type HadithCollection } from '~/utils/hadith-collection';

export function HadithSettings() {
  const location = useLocation();
  const navigate = useNavigate();
  return <details className='relative'>
    <summary className='list-none cursor-pointer p-1 text-gray-500 hover:text-black' title='Hadith settings' aria-label='Hadith settings'><Settings2 className='size-4' /></summary>
    <div className='absolute left-0 top-8 z-50 w-56 rounded-lg border border-gray-200 bg-white p-3 shadow-lg'>
      <label className='flex flex-col gap-2 text-sm'>Collection
        <select className='w-full rounded border border-gray-300 p-2' value={collectionFromPath(location.pathname)} onChange={event => {
          const collection = event.target.value as HadithCollection;
          event.target.closest('details')?.removeAttribute('open');
          navigate(hadithPath(collection));
        }}>
          <option value='bukhari'>Sahih Bukhari</option><option value='muslim'>Sahih Muslim</option>
        </select>
      </label>
    </div>
  </details>;
}
