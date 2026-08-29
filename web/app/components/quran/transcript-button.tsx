import { ClipboardCopy, FileText, Share2, X } from 'lucide-react';
import { useMemo, useState } from 'react';

interface Verse {
  number: number;
  text: string;
  arabic: string;
}

interface TranscriptButtonProps {
  chapterNumber: number;
  chapterName: string;
  chapterEnglish: string;
  verses: Verse[];
}

type TranscriptFormat = 'english' | 'arabic' | 'both';

export function TranscriptButton({
  chapterNumber,
  chapterName,
  chapterEnglish,
  verses,
}: TranscriptButtonProps) {
  const [open, setOpen] = useState(false);
  const [copied, setCopied] = useState(false);

  const nonBismillah = verses.filter((v) => v.number !== 0);
  const maxVerse = nonBismillah.length > 0 ? nonBismillah[nonBismillah.length - 1].number : 1;

  const [fromVerse, setFromVerse] = useState(1);
  const [toVerse, setToVerse] = useState(Math.min(10, maxVerse));
  const [format, setFormat] = useState<TranscriptFormat>('english');

  const transcript = useMemo(() => {
    const selected = nonBismillah.filter(
      (v) => v.number >= fromVerse && v.number <= toVerse
    );

    const header = `${chapterEnglish} (${chapterName}) — ${chapterNumber}:${fromVerse}${fromVerse !== toVerse ? `-${toVerse}` : ''}`;
    const lines = selected.map((v) => {
      const ref = `[${chapterNumber}:${v.number}]`;
      if (format === 'arabic') return `${v.arabic} ${ref}`;
      if (format === 'english') return `${v.text} ${ref}`;
      return `${v.arabic}\n${v.text} ${ref}`;
    });

    return `${header}\n\n${lines.join('\n\n')}`;
  }, [nonBismillah, fromVerse, toVerse, format, chapterNumber, chapterName, chapterEnglish]);

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(transcript);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  };

  const handleShare = async () => {
    if (navigator.share) {
      try {
        await navigator.share({
          title: `${chapterEnglish} ${chapterNumber}:${fromVerse}-${toVerse}`,
          text: transcript,
        });
      } catch {}
    } else {
      handleCopy();
    }
  };

  if (!open) {
    return (
      <button
        onClick={() => setOpen(true)}
        className='inline-flex items-center justify-center p-1 rounded-md transition-colors hover:bg-gray-100'
        title='Transcript'
        aria-label='Transcript'
      >
        <FileText className='size-4 text-gray-400 hover:text-gray-600' />
      </button>
    );
  }

  return (
    <div className='fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/40' onClick={() => setOpen(false)}>
      <div
        className='bg-white w-full sm:max-w-md sm:rounded-xl rounded-t-xl p-4 sm:p-5 max-h-[80vh] flex flex-col'
        onClick={(e) => e.stopPropagation()}
      >
        <div className='flex items-center justify-between mb-4'>
          <h3 className='font-semibold text-base sm:text-lg'>Copy Transcript</h3>
          <button onClick={() => setOpen(false)} className='p-1 hover:bg-gray-100 rounded'>
            <X className='size-4' />
          </button>
        </div>

        {/* Range picker */}
        <div className='flex items-center gap-2 mb-3'>
          <label className='text-sm text-gray-600'>Verses</label>
          <input
            type='number'
            min={1}
            max={toVerse}
            value={fromVerse}
            onChange={(e) => setFromVerse(Math.max(1, Math.min(Number(e.target.value), toVerse)))}
            className='w-16 border border-gray-300 rounded px-2 py-1 text-sm'
          />
          <span className='text-sm text-gray-400'>to</span>
          <input
            type='number'
            min={fromVerse}
            max={maxVerse}
            value={toVerse}
            onChange={(e) => setToVerse(Math.max(fromVerse, Math.min(Number(e.target.value), maxVerse)))}
            className='w-16 border border-gray-300 rounded px-2 py-1 text-sm'
          />
          <span className='text-xs text-gray-400'>of {maxVerse}</span>
        </div>

        {/* Format picker */}
        <div className='flex gap-2 mb-4'>
          {(['english', 'arabic', 'both'] as const).map((f) => (
            <button
              key={f}
              onClick={() => setFormat(f)}
              className={`text-xs px-2.5 py-1 rounded-md border capitalize transition-colors ${
                format === f
                  ? 'bg-black text-white border-black'
                  : 'border-gray-300 text-gray-600 hover:bg-gray-100'
              }`}
            >
              {f}
            </button>
          ))}
        </div>

        {/* Preview */}
        <div className='flex-1 overflow-y-auto bg-gray-50 rounded-lg p-3 mb-4 text-sm leading-relaxed whitespace-pre-wrap max-h-60 border border-gray-100'>
          {transcript}
        </div>

        {/* Actions */}
        <div className='flex gap-2'>
          <button
            onClick={handleCopy}
            className='flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg bg-black text-white text-sm font-medium hover:bg-gray-800 transition-colors'
          >
            <ClipboardCopy className='size-3.5' />
            {copied ? 'Copied!' : 'Copy'}
          </button>
          <button
            onClick={handleShare}
            className='flex-1 flex items-center justify-center gap-1.5 py-2 rounded-lg border border-gray-300 text-sm font-medium hover:bg-gray-100 transition-colors'
          >
            <Share2 className='size-3.5' />
            Share
          </button>
        </div>
      </div>
    </div>
  );
}
