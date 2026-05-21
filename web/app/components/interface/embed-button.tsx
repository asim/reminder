import { Code2 } from 'lucide-react';
import { useState } from 'react';
import { cn } from '~/utils/classname';

interface EmbedButtonProps {
  path: string;
  className?: string;
}

export function EmbedButton({ path, className }: EmbedButtonProps) {
  const [showCode, setShowCode] = useState(false);
  const [copied, setCopied] = useState(false);

  const embedUrl = `https://reminder.dev/embed${path}`;
  const iframeCode = `<iframe src="${embedUrl}" width="100%" height="220" style="border:1px solid #e5e7eb;border-radius:8px;" frameborder="0"></iframe>`;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(iframeCode);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    } catch {}
  };

  return (
    <>
      <button
        onClick={() => setShowCode(!showCode)}
        className={cn(
          'inline-flex items-center justify-center p-1 rounded-md transition-colors hover:bg-gray-100',
          className
        )}
        title='Embed'
        aria-label='Embed'
      >
        <Code2
          className={cn(
            'size-4',
            showCode ? 'text-blue-500' : 'text-gray-400 hover:text-gray-600'
          )}
        />
      </button>

      {showCode && (
        <div className='fixed inset-0 z-50 flex items-end sm:items-center justify-center bg-black/40' onClick={() => setShowCode(false)}>
          <div
            className='bg-white w-full sm:max-w-md sm:rounded-xl rounded-t-xl p-4 sm:p-5'
            onClick={(e) => e.stopPropagation()}
          >
            <div className='text-sm font-medium mb-2'>Embed this content</div>
            <p className='text-xs text-gray-500 mb-3'>
              Copy the code below and paste it into your website or blog.
            </p>
            <pre className='bg-gray-50 border border-gray-200 rounded-lg p-3 text-xs overflow-x-auto whitespace-pre-wrap break-all mb-3'>
              {iframeCode}
            </pre>
            <div className='flex gap-2'>
              <button
                onClick={handleCopy}
                className='flex-1 py-2 rounded-lg bg-black text-white text-sm font-medium hover:bg-gray-800 transition-colors'
              >
                {copied ? 'Copied!' : 'Copy code'}
              </button>
              <button
                onClick={() => setShowCode(false)}
                className='flex-1 py-2 rounded-lg border border-gray-300 text-sm font-medium hover:bg-gray-100 transition-colors'
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </>
  );
}
