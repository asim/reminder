import { Settings, X } from 'lucide-react';
import { createContext, useContext, useRef, useState, useId, useEffect, useLayoutEffect } from 'react';
import type { CSSProperties, ReactNode } from 'react';
import { useLocation, useNavigate } from 'react-router';
import type { ViewMode } from './view-mode';

export type Translation = 'khattab' | 'asad' | 'sahih';
type Preferences = {
  mode: ViewMode;
  translation: Translation;
  wordByWord: boolean;
  transliteration: boolean;
  commentary: boolean;
  arabicSize: number;
  englishSize: number;
};
const defaults: Preferences = {
  mode: 'translation', translation: 'khattab', wordByWord: false,
  transliteration: false, commentary: false, arabicSize: 100, englishSize: 100,
};
const storageKey = 'quran-reading-settings';
const validMode = (v: unknown): v is ViewMode => ['translation', 'arabic', 'english'].includes(String(v));
const validTranslation = (v: unknown): v is Translation => v === 'khattab' || v === 'asad' || v === 'sahih';
const validSize = (v: unknown) => typeof v === 'number' && [80, 100, 120, 140, 160].includes(v);
function savedPreferences(): Preferences {
  if (typeof window === 'undefined') return defaults;
  try {
    const saved = JSON.parse(localStorage.getItem(storageKey) || '{}');
    const oldMode = localStorage.getItem('quran-view-mode');
    return {
      mode: validMode(saved.mode) ? saved.mode : validMode(oldMode) ? oldMode : defaults.mode,
      translation: validTranslation(saved.translation) ? saved.translation : defaults.translation,
      wordByWord: typeof saved.wordByWord === 'boolean' ? saved.wordByWord : localStorage.getItem('quran-word-by-word-toggle') === 'true',
      transliteration: typeof saved.transliteration === 'boolean' ? saved.transliteration : localStorage.getItem('quran-transliteration-toggle') === 'true',
      commentary: typeof saved.commentary === 'boolean' ? saved.commentary : false,
      arabicSize: validSize(saved.arabicSize) ? saved.arabicSize : 100,
      englishSize: validSize(saved.englishSize) ? saved.englishSize : 100,
    };
  } catch { return defaults; }
}
const Context = createContext<{
  settings: Preferences;
  update: (patch: Partial<Preferences>) => void;
  textStyle: CSSProperties;
} | null>(null);

export function QuranSettingsProvider({ children }: { children: ReactNode }) {
  const [saved, setSaved] = useState(savedPreferences);
  const [pending, setPending] = useState<{ settings: Preferences; search: string } | null>(null);
  const readingPosition = useRef<{ element: Element; offset: number; scroller: HTMLElement } | null>(null);
  const location = useLocation();
  const navigate = useNavigate();
  const params = new URLSearchParams(location.search);
  const mode = params.get('mode');
  const translation = params.get('translation');
  const booleanParam = (key: string, fallback: boolean) => {
    const value = params.get(key);
    return value === null ? fallback : value === '1' || value === 'true';
  };
  const settings: Preferences = pending?.settings ?? {
    ...saved,
    mode: validMode(mode) ? mode : saved.mode,
    translation: validTranslation(translation) ? translation : saved.translation,
    wordByWord: booleanParam('wbw', saved.wordByWord),
    transliteration: booleanParam('transliteration', saved.transliteration),
    commentary: booleanParam('commentary', saved.commentary),
  };
  useEffect(() => {
    if (pending && location.search === '?' + pending.search) setPending(null);
  }, [location.search, pending]);
  useLayoutEffect(() => {
    const position = readingPosition.current;
    if (!position) return;
    const frame = requestAnimationFrame(() => {
      const { element, offset, scroller } = position;
      if (element.isConnected) scroller.scrollTop += element.getBoundingClientRect().top - scroller.getBoundingClientRect().top - offset;
      readingPosition.current = null;
    });
    return () => cancelAnimationFrame(frame);
  }, [location.key]);
  function update(patch: Partial<Preferences>) {
    const scroller = document.querySelector<HTMLElement>('[data-quran-scroll]');
    if (scroller) {
      const top = scroller.getBoundingClientRect().top;
      const element = Array.from(scroller.querySelectorAll('[data-chapter-verse]')).find(el => el.getBoundingClientRect().bottom > top);
      if (element) readingPosition.current = { element, scroller, offset: element.getBoundingClientRect().top - top };
    }
    const next = { ...settings, ...patch };
    setSaved(next);
    try { localStorage.setItem(storageKey, JSON.stringify(next)); } catch { /* Still works without storage. */ }
    params.set('mode', next.mode);
    params.set('translation', next.translation);
    params.set('wbw', next.wordByWord ? '1' : '0');
    params.set('transliteration', next.transliteration ? '1' : '0');
    params.set('commentary', next.commentary ? '1' : '0');
    setPending({ settings: next, search: params.toString() });
    navigate({ pathname: location.pathname, search: params.toString(), hash: location.hash }, { replace: true, preventScrollReset: true });
  }
  const textStyle = {
    '--quran-arabic-scale': settings.arabicSize / 100,
    '--quran-english-scale': settings.englishSize / 100,
  } as CSSProperties;
  return <Context.Provider value={{ settings, update, textStyle }}>{children}</Context.Provider>;
}

export function useQuranSettings() {
  const value = useContext(Context);
  if (!value) throw new Error('Missing QuranSettingsProvider');
  return value;
}

export function ReadingSettings() {
  const { settings, update } = useQuranSettings();
  const location = useLocation();
  const dialog = useRef<HTMLDialogElement>(null);
  const titleId = useId();
  if (!/^\/quran\/\d+(?:\/\d+)?\/?$/.test(location.pathname)) return null;
  const selectClass = 'w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm';
  return <>
    <button type='button' onClick={() => dialog.current?.showModal()} aria-label='Reading settings' title='Reading settings'
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
        <h2 id={titleId} className='text-lg font-semibold'>Reading settings</h2>
        <button type='button' aria-label='Close reading settings' onClick={() => dialog.current?.close()} className='p-2 rounded hover:bg-gray-100'><X className='size-5' /></button>
      </div>
      <div className='space-y-5'>
        <label className='block text-sm'>
          <span className='block mb-2 font-medium'>Display</span>
          <select aria-label='Display' className={selectClass} value={settings.mode} onChange={e => update({ mode: e.target.value as ViewMode })}>
            <option value='translation'>Arabic and English</option>
            <option value='arabic'>Arabic only</option>
            <option value='english'>English only</option>
          </select>
        </label>
        <label className='block text-sm'>
          <span className='block mb-2 font-medium'>English translation</span>
          <select aria-label='English translation' className={selectClass} value={settings.translation} onChange={e => update({ translation: e.target.value as Translation })}>
            <option value='khattab'>Mustafa Khattab · The Clear Quran</option>
            <option value='sahih'>Saheeh International</option>
            <option value='asad'>Muhammad Asad · The Message of the Qur’an</option>
          </select>
        </label>
        <div className='grid grid-cols-2 gap-4'>
          {(['arabicSize', 'englishSize'] as const).map(key => <label key={key} className='text-sm'>
            <span className='block mb-2 font-medium'>{key === 'arabicSize' ? 'Arabic text size' : 'English text size'}</span>
            <select aria-label={key === 'arabicSize' ? 'Arabic text size' : 'English text size'} className={selectClass} value={settings[key]} onChange={e => update({ [key]: Number(e.target.value) })}>
              {[80, 100, 120, 140, 160].map(size => <option key={size} value={size}>{size === 100 ? 'Default' : `${size}%`}</option>)}
            </select>
          </label>)}
        </div>
        <fieldset className='border-t border-gray-100 pt-3'>
          <legend className='sr-only'>Reading aids</legend>
          {([
            ['wordByWord', 'Word-by-word'],
            ['transliteration', 'Transliteration'],
            ['commentary', 'Commentary'],
          ] as const).map(([key, label]) => <label key={key} className='flex items-center justify-between gap-4 py-3 text-sm cursor-pointer'>
            <span>{label}</span>
            <input type='checkbox' checked={settings[key]} onChange={e => update({ [key]: e.target.checked })} className='accent-black h-4 w-4' />
          </label>)}
        </fieldset>
        <p className='text-xs text-gray-500 leading-relaxed'>
          Translation applies to verse text. Word meanings and stored tafsir are separate sources. English audio uses Sahih International.
          {settings.translation === 'asad' && <> Asad’s explanatory notes are available <a className='underline' href={`https://www.alim.org/translation/asad/${location.pathname.split('/')[2]}/`} target='_blank' rel='noreferrer'>on Alim</a>.</>}
        </p>
        <button type='button' onClick={() => update(defaults)} className='text-sm underline underline-offset-4'>Reset reading settings</button>
      </div>
    </dialog>
  </>;
}
