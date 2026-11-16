import { useQuranSettings } from './reading-settings';

import { tafsirEditions, type TafsirEdition } from '~/utils/tafsir-editions';

export function TafsirSelect() {
  const { settings, update } = useQuranSettings();
  return <label className='block text-left text-sm'>
    <span className='mb-2 block font-medium'>Tafsir source</span>
    <select className='w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm'
      value={settings.tafsir} onChange={event => update({ tafsir: event.target.value as TafsirEdition })}>
      {tafsirEditions.map(edition => <option key={edition.id} value={edition.id}>{edition.label}</option>)}
    </select>
    {settings.tafsir === 'original-en' && <span className='mt-2 block text-xs leading-relaxed text-gray-500'>
      Preserved from Reminder’s original dataset. Author and source have not been verified.
    </span>}
  </label>;
}
