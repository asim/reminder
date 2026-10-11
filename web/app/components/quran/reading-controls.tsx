import { ViewMode } from './view-mode';

export function ReadingControls({ mode, setMode, wordByWord, setWordByWord, transliteration, setTransliteration, commentary, setCommentary }: {
  mode: import('./view-mode').ViewMode; setMode: (mode: import('./view-mode').ViewMode) => void;
  wordByWord: boolean; setWordByWord: (value: boolean) => void;
  transliteration: boolean; setTransliteration: (value: boolean) => void;
  commentary: boolean; setCommentary: (value: boolean) => void;
}) {
  // Preserve the visible verse when annotations above it change the layout.
  function preservePosition(change: () => void) {
    const anchor = Array.from(document.querySelectorAll<HTMLElement>('[data-chapter-verse]'))
      .find(el => el.getBoundingClientRect().bottom > 160);
    const top = anchor?.getBoundingClientRect().top;
    change();
    requestAnimationFrame(() => {
      const current = anchor && document.querySelector<HTMLElement>(`[data-chapter-verse="${anchor.dataset.chapterVerse}"]`);
      if (current && top !== undefined) current.closest('.overflow-y-auto')?.scrollBy(0, current.getBoundingClientRect().top - top);
    });
  }
  return (
    <div className='sticky top-0 z-20 bg-white border-b border-gray-200 py-2 mb-4' aria-label='Reading options'>
      <ViewMode mode={mode} onChange={value => preservePosition(() => setMode(value))} />
      <div className='flex flex-wrap gap-x-4 gap-y-2 px-1 text-sm'>
        {mode !== 'english' && <>
          <label className='flex items-center gap-2'><input type='checkbox' checked={wordByWord} onChange={e => preservePosition(() => setWordByWord(e.target.checked))} />Word-by-word</label>
          <label className='flex items-center gap-2'><input type='checkbox' checked={transliteration} onChange={e => preservePosition(() => setTransliteration(e.target.checked))} />Transliteration</label>
        </>}
        <label className='flex items-center gap-2'><input type='checkbox' checked={commentary} onChange={e => preservePosition(() => setCommentary(e.target.checked))} />Tafsir</label>
      </div>
    </div>
  );
}
