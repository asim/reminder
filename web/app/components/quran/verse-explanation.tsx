import { useMutation } from '@tanstack/react-query';
import { Link } from 'react-router';
import { httpPost } from '~/utils/http';

type Explanation = {
  text: string;
  generated: boolean;
  source: string;
  references: { text: string; metadata: Record<string, string> }[];
};

export function VerseExplanation({ chapter, verse }: { chapter: number; verse: number }) {
  const explain = useMutation({
    mutationFn: () => httpPost<Explanation>(`/api/quran/${chapter}/${verse}/explain`, {}),
    retry: false,
  });
  return <section className='my-4 space-y-3'>
    <div className='flex flex-wrap items-center gap-3'>
      <button type='button' className='rounded-md border border-gray-300 px-3 py-2 text-sm hover:bg-gray-50 disabled:opacity-60' disabled={explain.isPending} onClick={() => explain.mutate()}>
        {explain.isPending ? 'Explaining…' : 'Explain verse'}
      </button>
      <Link className='text-sm text-blue-700 underline underline-offset-4' to={`/search?mode=ask&q=${encodeURIComponent(`What can I learn from Quran ${chapter}:${verse}?`)}`}>Ask about this verse</Link>
    </div>
    <div aria-live='polite'>
      {explain.error && <p role='alert' className='text-sm'>{explain.error.message}</p>}
      {explain.data && <div className='space-y-3 rounded-lg border border-gray-200 p-4'>
        <h2 className='font-semibold'>{explain.data.generated ? 'Explanation' : explain.data.source || 'Stored commentary'}</h2>
        <p className='whitespace-pre-wrap text-sm sm:text-base leading-relaxed'>{explain.data.text}</p>
        {explain.data.references.length > 0 && <details>
          <summary className='cursor-pointer text-sm font-medium'>Sources ({explain.data.references.length})</summary>
          <ol className='mt-3 space-y-3 text-sm'>
            {explain.data.references.map((ref, i) => <li key={i}>
              <p className='font-medium'>[{i + 1}] {ref.metadata.title || ref.metadata.translator || 'Quran'} · {ref.metadata.ayahs || `${chapter}:${verse}`}</p>
              <p className='mt-1 whitespace-pre-wrap'>{ref.text}</p>
            </li>)}
          </ol>
        </details>}
      </div>}
    </div>
  </section>;
}
