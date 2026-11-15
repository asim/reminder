import { useModelAvailable } from '~/queries/capabilities';
import { useEffect, useRef, useState } from 'react';
import { BookOpen, X } from 'lucide-react';

type Reference = { text: string; metadata: Record<string, string> };
type Turn = { question: string; answer: string; references: Reference[]; coverage?: { included_verses: number; total_verses: number } };

type Commentary = { number: number; comments?: string; commentary_source?: string };

export function ExplainButton({ chapter, verse, translation, commentary = [] }: { chapter: number; verse?: number; translation: string; commentary?: Commentary[] }) {
  const modelAvailable = useModelAvailable();
  const passages = commentary.filter(item => item.number > 0 && !!item.comments);
  const [view, setView] = useState<'explain' | 'tafsir'>('tafsir');
  const title = view === 'explain' ? 'Explain' : 'Tafsir';
  const dialog = useRef<HTMLDialogElement>(null);
  const pending = useRef<AbortController | null>(null);
  const [opened, setOpened] = useState(false);
  const [turns, setTurns] = useState<Turn[]>([]);
  const [question, setQuestion] = useState('');
  const [selection, setSelection] = useState('');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const label = verse ? `Quran ${chapter}:${verse}` : `Quran chapter ${chapter}`;

  useEffect(() => {
    pending.current?.abort(); dialog.current?.close();
    setOpened(false); setTurns([]); setQuestion(''); setSelection(''); setError(''); setBusy(false);
    return () => { pending.current?.abort(); };
  }, [chapter, verse, translation]);

  useEffect(() => { if (opened) dialog.current?.showModal(); }, [opened]);

  useEffect(() => {
    if (!modelAvailable) {
      pending.current?.abort(); pending.current = null;
      setBusy(false); setView('tafsir');
    }
  }, [modelAvailable]);

  function close() {
    pending.current?.abort(); pending.current = null;
    setBusy(false); dialog.current?.close(); setOpened(false);
  }

  function open() {
    const selected = window.getSelection();
    const element = selected?.anchorNode?.parentElement?.closest('[data-chapter-verse]');
    const end = selected?.focusNode?.parentElement?.closest('[data-chapter-verse]');
    const reference = element?.getAttribute('data-chapter-verse');
    const valid = element && element === end && (verse ? reference === `${chapter}:${verse}` : reference?.startsWith(`${chapter}:`));
    const text = valid ? selected?.toString().trim().slice(0, 2000) || '' : '';
    if (text !== selection) { setTurns([]); setQuestion(''); }
    setSelection(text); setError(''); setView(modelAvailable ? 'explain' : 'tafsir'); setOpened(true);
  }

  async function ask() {
    if (busy || !modelAvailable) return;
    const prompt = question.trim() || 'Explain this passage';
    const controller = new AbortController(); pending.current = controller;
    setBusy(true); setError('');
    try {
      const response = await fetch('/api/explain', {
        method: 'POST', headers: { 'Content-Type': 'application/json' }, signal: controller.signal,
        body: JSON.stringify({ context: { source: 'quran', chapter, verse, translation, selection }, question: question.trim(),
          history: turns.slice(-4).flatMap(turn => [{ role: 'user', content: turn.question }, { role: 'assistant', content: turn.answer }]) }),
      });
      const result = await response.json();
      if (!response.ok) throw new Error(result.message || 'Explanation unavailable. Please try again.');
      if (pending.current !== controller) return;
      setTurns(previous => [...previous, { question: prompt, answer: result.answer, references: result.references || [], coverage: result.coverage }]); setQuestion('');
    } catch (err) {
      if (!controller.signal.aborted) setError(err instanceof Error ? err.message : 'Unable to load explanation.');
    } finally {
      if (pending.current === controller) { pending.current = null; setBusy(false); }
    }
  }

  if (!modelAvailable && !passages.length) return null;

  return <>
    <button type='button' onClick={open} className='inline-flex items-center gap-1 rounded-md p-1 text-xs text-gray-500 hover:bg-gray-100' aria-label={`${modelAvailable ? 'Explain' : 'Tafsir'} ${label}`}>
      <BookOpen className='size-4' /> {modelAvailable ? 'Explain' : 'Tafsir'}
    </button>
    {opened && <dialog ref={dialog} aria-label={`${title} ${label}`} onCancel={close} className='fixed inset-x-0 bottom-0 top-auto m-0 h-[92dvh] max-h-[92dvh] w-full max-w-none overflow-y-auto overscroll-contain rounded-t-xl border border-gray-200 bg-white p-0 pb-[env(safe-area-inset-bottom)] text-gray-900 shadow-xl backdrop:bg-black/30 md:inset-y-0 md:left-auto md:right-0 md:h-dvh md:max-h-dvh md:w-[32rem] md:rounded-none'>
      <div className='sticky top-0 flex items-center justify-between gap-4 border-b bg-white p-4'>
        <h2 className='font-medium'>{title} {label}</h2>
        <button type='button' onClick={close} aria-label='Close reading panel' className='rounded p-2 hover:bg-gray-100'><X className='size-4' /></button>
      </div>
      {modelAvailable && passages.length > 0 && <nav aria-label='Reading panel' className='flex gap-4 border-b px-4 py-2 text-sm'>
        <button type='button' aria-pressed={view === 'explain'} onClick={() => setView('explain')} className={view === 'explain' ? 'font-medium underline underline-offset-4' : 'text-gray-500'}>Explain</button>
        <button type='button' aria-pressed={view === 'tafsir'} onClick={() => setView('tafsir')} className={view === 'tafsir' ? 'font-medium underline underline-offset-4' : 'text-gray-500'}>Tafsir</button>
      </nav>}
      {view === 'tafsir' ? <div className='space-y-6 p-4 text-sm'>
        {passages.map(item => <section key={item.number} className='space-y-2'>
          <h3 className='font-medium'>{chapter}:{item.number} — {item.commentary_source || 'Stored commentary'}</h3>
          <p className='whitespace-pre-wrap leading-relaxed'>{item.comments}</p>
        </section>)}
        {!passages.length && <p>No stored tafsir is available for this passage.</p>}
      </div> : <div className='space-y-4 p-4 text-sm' aria-live='polite'>
        {selection && <blockquote className='border-l-2 pl-3' dir='auto'>{selection}</blockquote>}
        {!turns.length && <p className='text-gray-600'>Explore this passage through its Arabic, translations and tafsir. Ask a question or choose Explain.</p>}
        {turns.map((turn, index) => <section key={index} className='space-y-3 border-b pb-4'>
          <h3 className='font-medium'>{turn.question}</h3>
          {turn.coverage && <p className='text-xs text-gray-500'>Source text covers {turn.coverage.included_verses} of {turn.coverage.total_verses} verses{turn.coverage.included_verses < turn.coverage.total_verses ? ' — a partial chapter overview.' : '.'}</p>}
          <p className='whitespace-pre-wrap leading-relaxed'>{turn.answer}</p>
          <details><summary className='cursor-pointer text-gray-600'>Sources ({turn.references.length})</summary>
            <ol className='mt-2 space-y-3'>{turn.references.map((ref, i) => <li key={i} className='break-words text-xs leading-relaxed'>
              <strong>[{i + 1}] {ref.metadata.title || ref.metadata.translator || ref.metadata.source}{ref.metadata.chapter && ` — ${ref.metadata.chapter}:${ref.metadata.verse}`}</strong>
              <p className='whitespace-pre-wrap'>{ref.text}</p>
            </li>)}</ol>
          </details>
        </section>)}
        <p className='text-xs text-gray-500'>AI-generated explanation. Check the cited sources.</p>
        {error && <p role='alert' className='text-red-700'>{error}{passages.length > 0 && ' You can read Tafsir above, or retry your question below.'}</p>}
        <form onSubmit={event => { event.preventDefault(); void ask(); }} className='space-y-2'>
          <label className='block' htmlFor={`explain-question-${chapter}-${verse || 0}`}>{turns.length ? 'Ask about this passage' : 'Your question (optional)'}</label>
          <textarea id={`explain-question-${chapter}-${verse || 0}`} value={question} onChange={event => setQuestion(event.target.value)} maxLength={2000} rows={2} className='w-full rounded border border-gray-300 p-2' />
          <button disabled={busy || (turns.length > 0 && !question.trim())} className='rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50'>{busy ? 'Explaining…' : turns.length ? 'Ask' : 'Explain'}</button>
        </form>
      </div>}
    </dialog>}
  </>;
}
