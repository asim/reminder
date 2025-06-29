import { useState, useEffect } from 'react';

const STORAGE_KEY = 'quran-word-by-word-toggle';

export function useWordByWordToggle(): [boolean, (val: boolean) => void] {
  const [enabled, setEnabled] = useState<boolean>(() => {
    if (typeof window !== 'undefined') {
      const stored = window.localStorage.getItem(STORAGE_KEY);
      return stored === 'true';
    }
    return false;
  });

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEY, String(enabled));
    }
  }, [enabled]);

  return [enabled, setEnabled];
}
