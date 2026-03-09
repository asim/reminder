import React, { useCallback, useEffect, useRef, useState } from 'react';

type ClickableArabicWordProps = {
  arabic: string;
  transliteration?: string;
  children?: React.ReactNode;
  className?: string;
};

export function ClickableArabicWord({
  arabic,
  transliteration,
  children,
  className = '',
}: ClickableArabicWordProps) {
  const [showTooltip, setShowTooltip] = useState(false);
  const ref = useRef<HTMLSpanElement>(null);

  const handleClick = useCallback(() => {
    if (!transliteration) return;
    setShowTooltip((prev) => !prev);
  }, [transliteration]);

  useEffect(() => {
    if (!showTooltip) return;

    const handleClickOutside = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setShowTooltip(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [showTooltip]);

  return (
    <span
      ref={ref}
      className={`verse-arabic-word ${className}`}
      onClick={handleClick}
      style={{ position: 'relative' }}
    >
      {children ?? arabic}
      {showTooltip && transliteration && (
        <span className='transliteration-tooltip'>{transliteration}</span>
      )}
    </span>
  );
}
