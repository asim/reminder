import { Settings, X } from 'lucide-react';
import { useId, useRef } from 'react';
import { useLocation, useNavigate } from 'react-router';
import { collectionFromPath, hadithPath, type HadithCollection } from '~/utils/hadith-collection';

export function HadithSettings() {
  const location = useLocation();
  const navigate = useNavigate();
  const dialog = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  const selectClass = 'w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm';
  return <>
    <button type='button' onClick={() => dialog.current?.showModal()} aria-label='Hadith settings' title='Hadith settings'
      className='relative before:absolute before:-inset-2 shrink-0 cursor-pointer rounded text-gray-500 hover:text-black focus-visible:outline-2 focus-visible:outline-black'>
      <Settings className='size-4' />
    </button>
    <dialog ref={dialog} aria-labelledby={titleId}
      className='reading-settings-dialog text-gray-900 bg-white p-5 shadow-xl backdrop:bg-black/40'
      onClick={event => { if (event.target === dialog.current) {
        const rect = dialog.current.getBoundingClientRect();
        if (event.clientX < rect.left || event.clientX > rect.right || event.clientY < rect.top || event.clientY > rect.bottom) dialog.current.close();
      } }}>
      <div className='flex items-center justify-between mb-5'>
        <h2 id={titleId} className='text-lg font-semibold'>Hadith settings</h2>
        <button type='button' aria-label='Close Hadith settings' onClick={() => dialog.current?.close()} className='p-2 rounded hover:bg-gray-100'><X className='size-5' /></button>
      </div>
      <div className='space-y-5'>
        <label className='block text-sm'>
          <span className='block mb-2 font-medium'>Collection</span>
          <select aria-label='Collection' className={selectClass} value={collectionFromPath(location.pathname)} onChange={event => {
            const collection = event.target.value as HadithCollection;
            dialog.current?.close();
            navigate(hadithPath(collection));
          }}>
            <option value='bukhari'>Sahih Bukhari</option>
            <option value='muslim'>Sahih Muslim</option>
          </select>
        </label>
      </div>
    </dialog>
  </>;
}
