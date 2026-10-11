import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router';

const STORAGE_KEY = 'quran-transliteration-toggle';

export function useTransliterationToggle(): [boolean, (val: boolean) => void] {
  const [searchParams] = useSearchParams();
  const [enabled, setEnabled] = useState<boolean>(() => {
    if (typeof window !== 'undefined') {
      const param = searchParams.get('transliteration');
      if (param !== null) return param === '1' || param === 'true';
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
