import { Star } from 'lucide-react';
import { useBookmarks, type BookmarkType } from '~/hooks/use-bookmarks';
import { cn } from '~/utils/classname';

interface BookmarkButtonProps {
  type: BookmarkType;
  itemKey: string;
  label: string;
  url: string;
  className?: string;
}

export function BookmarkButton({
  type,
  itemKey,
  label,
  url,
  className,
}: BookmarkButtonProps) {
  const { hasBookmark, toggleBookmark } = useBookmarks();
  const isBookmarked = hasBookmark(type, itemKey);

  const handleClick = () => {
    toggleBookmark(type, itemKey, label, url);
  };

  return (
    <button
      onClick={handleClick}
      className={cn(
        'inline-flex items-center justify-center p-1 rounded-md transition-colors hover:bg-gray-100',
        className
      )}
      title={isBookmarked ? 'Remove bookmark' : 'Add bookmark'}
      aria-label={isBookmarked ? 'Remove bookmark' : 'Add bookmark'}
    >
      <Star
        className={cn(
          'size-4',
          isBookmarked
            ? 'fill-yellow-400 text-yellow-400'
            : 'text-gray-400 hover:text-gray-600'
        )}
      />
    </button>
  );
}
