import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router';

const STORAGE_KEY = 'quran-word-by-word-toggle';

export function useWordByWordToggle(): [boolean, (val: boolean) => void] {
  const [searchParams] = useSearchParams();

  const [enabled, setEnabled] = useState<boolean>(() => {
    if (typeof window === 'undefined') return false;

    const param = searchParams.get('wbw');
    if (param !== null) return param === '1' || param === 'true';

    return window.localStorage.getItem(STORAGE_KEY) === 'true';
  });

  useEffect(() => {
    if (typeof window !== 'undefined') {
      window.localStorage.setItem(STORAGE_KEY, String(enabled));
    }
  }, [enabled]);

  return [enabled, setEnabled];
}
