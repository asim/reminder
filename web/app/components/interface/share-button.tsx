import { Share2 } from 'lucide-react';
import { useState } from 'react';
import { cn } from '~/utils/classname';

interface ShareButtonProps {
  title: string;
  text?: string;
  url: string;
  className?: string;
}

export function ShareButton({ title, text, url, className }: ShareButtonProps) {
  const [copied, setCopied] = useState(false);

  const handleClick = async () => {
    // Resolve relative URLs to absolute
    const absoluteUrl = url.startsWith('http')
      ? url
      : typeof window !== 'undefined'
        ? new URL(url, window.location.origin).toString()
        : url;

    const shareData: ShareData = {
      title,
      url: absoluteUrl,
    };
    if (text) shareData.text = text;

    // Prefer the Web Share API (native share sheet on mobile/desktop)
    if (typeof navigator !== 'undefined' && navigator.share) {
      try {
        await navigator.share(shareData);
        return;
      } catch (err) {
        // User cancelled or share failed — fall through to clipboard
        if ((err as Error).name === 'AbortError') return;
      }
    }

    // Fallback: copy link to clipboard
    if (typeof navigator !== 'undefined' && navigator.clipboard) {
      try {
        await navigator.clipboard.writeText(absoluteUrl);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
      } catch {
        // Clipboard blocked — nothing more we can do
      }
    }
  };

  return (
    <button
      onClick={handleClick}
      className={cn(
        'inline-flex items-center justify-center p-1 rounded-md transition-colors hover:bg-gray-100',
        className
      )}
      title={copied ? 'Link copied!' : 'Share'}
      aria-label='Share'
    >
      <Share2
        className={cn(
          'size-4',
          copied ? 'text-green-500' : 'text-gray-400 hover:text-gray-600'
        )}
      />
    </button>
  );
}
