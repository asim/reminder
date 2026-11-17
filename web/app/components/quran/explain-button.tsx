import { TafsirSelect } from './tafsir-select';
import Markdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { PrimaryButton } from '~/components/interface/primary-button';
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



  return <>
    <button type='button' onClick={open} className='inline-flex items-center justify-center p-1 rounded-md transition-colors text-gray-400 hover:bg-gray-100 hover:text-gray-600' title={`${modelAvailable ? 'Explain' : 'Tafsir'} ${label}`} aria-label={`${modelAvailable ? 'Explain' : 'Tafsir'} ${label}`}>
      <BookOpen className='size-4' aria-hidden='true' />
    </button>
    {opened && <dialog ref={dialog} aria-label={`${title} ${label}`} onCancel={close} className='reading-settings-dialog explanation-dialog bg-white text-left text-gray-900 shadow-xl backdrop:bg-black/40'>
      <div className='sticky top-0 z-10 flex items-center justify-between gap-4 border-b border-gray-100 bg-white p-4 sm:p-5'>
        <h2 className='text-lg font-semibold'>{title} {label}</h2>
        <button type='button' onClick={close} aria-label='Close reading panel' className='rounded p-2 hover:bg-gray-100'><X className='size-4' /></button>
      </div>
      {modelAvailable && <nav aria-label='Reading panel' className='flex gap-4 border-b border-gray-100 px-4 py-3 text-sm sm:px-5'>
        <button type='button' aria-pressed={view === 'explain'} onClick={() => setView('explain')} className={view === 'explain' ? 'font-medium underline underline-offset-4' : 'text-gray-500'}>Explain</button>
        <button type='button' aria-pressed={view === 'tafsir'} onClick={() => setView('tafsir')} className={view === 'tafsir' ? 'font-medium underline underline-offset-4' : 'text-gray-500'}>Tafsir</button>
      </nav>}
      {view === 'tafsir' ? <div className='space-y-6 p-4 text-left text-sm leading-relaxed sm:p-5'>
        <TafsirSelect />
        {passages.map(item => <section key={item.number} className='space-y-2'>
          <h3 className='font-medium'>{chapter}:{item.number} — {item.commentary_source || 'Stored commentary'}</h3>
          <p className='whitespace-pre-wrap leading-relaxed'>{item.comments}</p>
        </section>)}
        {!passages.length && <p>This tafsir has no saved commentary for this passage. Choose another source, or reconnect to update the chapter.</p>}
      </div> : <div className='space-y-4 p-4 text-left text-sm sm:p-5' aria-live='polite'>
        {selection && <div><blockquote className='border-l-2 pl-3' dir='auto'>{selection}</blockquote><button type='button' className='mt-1 text-xs underline' onClick={() => setSelection('')}>Clear selection</button></div>}
        {!turns.length && <p className='text-gray-600'>Explore this passage through its Arabic, translations and tafsir. Ask a question or choose Explain.</p>}
        {turns.map((turn, index) => <section key={index} className='space-y-3 border-b pb-4'>
          <h3 className='font-medium'>{turn.question}</h3>
          {turn.coverage && <p className='text-xs text-gray-500'>Source text covers {turn.coverage.included_verses} of {turn.coverage.total_verses} verses{turn.coverage.included_verses < turn.coverage.total_verses ? ' — a partial chapter overview.' : '.'}</p>}
          <div className='prose prose-sm max-w-none break-words text-left prose-headings:font-semibold prose-headings:text-gray-900 prose-p:leading-relaxed prose-pre:overflow-x-auto [&_table]:block [&_table]:overflow-x-auto'>
            <Markdown remarkPlugins={[remarkGfm]} skipHtml components={{ img: () => null }}>{turn.answer}</Markdown>
          </div>
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
          <textarea id={`explain-question-${chapter}-${verse || 0}`} value={question} onChange={event => setQuestion(event.target.value)} maxLength={2000} rows={2} className='w-full rounded-md border border-gray-300 bg-white px-3 py-2 text-sm' />
          <PrimaryButton disabled={busy || (turns.length > 0 && !question.trim())}>{busy ? 'Explaining…' : turns.length ? 'Ask' : 'Explain'}</PrimaryButton>
        </form>
      </div>}
    </dialog>}
  </>;
}
