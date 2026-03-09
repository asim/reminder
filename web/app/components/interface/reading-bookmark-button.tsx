import { BookmarkCheck } from 'lucide-react';
import { useReadingBookmark, type ReadingBookmarkType } from '~/hooks/use-reading-bookmark';
import { cn } from '~/utils/classname';

interface ReadingBookmarkButtonProps {
  type: ReadingBookmarkType;
  label: string;
  url: string;
  excerpt?: string;
  className?: string;
}

export function ReadingBookmarkButton({
  type,
  label,
  url,
  excerpt,
  className,
}: ReadingBookmarkButtonProps) {
  const { isReadingBookmark, setReadingBookmark, clearReadingBookmark } = useReadingBookmark();
  const isActive = isReadingBookmark(type, url);

  const handleClick = () => {
    if (isActive) {
      clearReadingBookmark(type);
    } else {
      setReadingBookmark(type, label, url, excerpt);
    }
  };

  return (
    <button
      onClick={handleClick}
      className={cn(
        'inline-flex items-center justify-center p-1 rounded-md transition-colors hover:bg-gray-100',
        className
      )}
      title={isActive ? 'Remove reading bookmark' : 'Set as reading bookmark'}
      aria-label={isActive ? 'Remove reading bookmark' : 'Set as reading bookmark'}
    >
      <BookmarkCheck
        className={cn(
          'size-4',
          isActive
            ? 'fill-blue-500 text-blue-500'
            : 'text-gray-400 hover:text-gray-600'
        )}
      />
    </button>
  );
}
