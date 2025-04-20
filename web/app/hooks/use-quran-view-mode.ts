import { useEffect, useState } from 'react';
import type { ViewMode } from '~/components/quran/view-mode';

const VIEW_MODE_KEY = 'quran-view-mode';

export function useQuranViewMode() {
  const [mode, setMode] = useState<ViewMode>(() => {
    if (typeof window === 'undefined') return 'translation';
    return (localStorage.getItem(VIEW_MODE_KEY) as ViewMode) || 'translation';
  });

  useEffect(() => {
    localStorage.setItem(VIEW_MODE_KEY, mode);
  }, [mode]);

  return [mode, setMode] as const;
} 