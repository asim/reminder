import { useEffect, useRef, useState } from 'react';
import { BookOpen, X } from 'lucide-react';

type Reference = { text: string; metadata: Record<string, string> };
type Turn = { question: string; answer: string; references: Reference[]; coverage?: { included_verses: number; total_verses: number } };

export function ExplainButton({ chapter, verse, translation }: { chapter: number; verse?: number; translation: string }) {
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
    setSelection(text); setError(''); setOpened(true);
  }

  async function ask() {
    if (busy) return;
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

  return <>
    <button type='button' onClick={open} className='inline-flex items-center gap-1 rounded-md p-1 text-xs text-gray-500 hover:bg-gray-100' aria-label={`Explain ${label}`}>
      <BookOpen className='size-4' /> Explain
    </button>
    {opened && <dialog ref={dialog} aria-label={`Explain ${label}`} onCancel={close} className='m-auto w-[calc(100%-1rem)] max-w-xl max-h-[85dvh] rounded-xl border border-gray-200 bg-white p-0 text-gray-900 shadow-xl backdrop:bg-black/40'>
      <div className='sticky top-0 flex items-center justify-between gap-4 border-b bg-white p-4'>
        <h2 className='font-medium'>Explain {label}</h2>
        <button type='button' onClick={close} aria-label='Close explanation' className='rounded p-2 hover:bg-gray-100'><X className='size-4' /></button>
      </div>
      <div className='space-y-4 p-4 text-sm' aria-live='polite'>
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
        {error && <p role='alert' className='text-red-700'>{error}</p>}
        <form onSubmit={event => { event.preventDefault(); void ask(); }} className='space-y-2'>
          <label className='block' htmlFor={`explain-question-${chapter}-${verse || 0}`}>{turns.length ? 'Ask about this passage' : 'Your question (optional)'}</label>
          <textarea id={`explain-question-${chapter}-${verse || 0}`} value={question} onChange={event => setQuestion(event.target.value)} maxLength={2000} rows={2} className='w-full rounded border border-gray-300 p-2' />
          <button disabled={busy || (turns.length > 0 && !question.trim())} className='rounded bg-gray-900 px-4 py-2 text-white disabled:opacity-50'>{busy ? 'Explaining…' : turns.length ? 'Ask' : 'Explain'}</button>
        </form>
      </div>
    </dialog>}
  </>;
}
