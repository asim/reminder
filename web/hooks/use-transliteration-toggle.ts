import { useEffect, useState } from 'react';

const STORAGE_KEY = 'quran-transliteration-toggle';

export function useTransliterationToggle(): [boolean, (val: boolean) => void] {
  const [enabled, setEnabled] = useState<boolean>(() => {
    if (typeof window !== 'undefined') {
      return window.localStorage.getItem(STORAGE_KEY) === 'true';
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
