import { useState } from 'react';

export function ExplainVerse({ chapter, verse, commentary }: { chapter: number; verse: number; commentary?: string }) {
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<{ text: string; generated: boolean }>();
  const [error, setError] = useState('');
  async function explain() {
    setOpen(!open);
    if (open || result || loading) return;
    setLoading(true); setError('');
    try {
      const response = await fetch(`/api/quran/${chapter}/${verse}/explain`, { method: 'POST' });
      if (!response.ok) throw new Error('Explanation could not be loaded. Please try again.');
      setResult(await response.json());
    } catch (e) {
      if (commentary) setResult({ text: commentary, generated: false });
      else setError(e instanceof Error ? e.message : 'Could not load explanation.');
    }
    finally { setLoading(false); }
  }
  return <div className='my-2 text-sm'>
    <button type='button' onClick={explain} aria-expanded={open} className='underline underline-offset-4'>{open ? 'Close explanation' : 'Explain'}</button>
    {open && <div className='mt-3 rounded border border-gray-200 bg-gray-50 p-3' aria-live='polite'>
      {loading && <p>Explaining…</p>}
      {error && <p role='alert'>{error}</p>}
      {result && <><p className='text-xs text-gray-500 mb-2'>{result.generated ? 'AI explanation based on this verse and stored commentary' : 'Stored commentary · AI explanation unavailable'} · Quran {chapter}:{verse}</p><div className='whitespace-pre-wrap leading-relaxed'>{result.text}</div></>}
    </div>}
  </div>;
}
