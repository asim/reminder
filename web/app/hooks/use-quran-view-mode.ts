import { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router';
import type { ViewMode } from '~/components/quran/view-mode';

const VIEW_MODE_KEY = 'quran-view-mode';
const VALID_MODES: ViewMode[] = ['translation', 'arabic', 'english'];

export function useQuranViewMode() {
  const [searchParams] = useSearchParams();

  const [mode, setMode] = useState<ViewMode>(() => {
    if (typeof window === 'undefined') return 'translation';

    const paramMode = searchParams.get('mode') as ViewMode | null;
    if (paramMode && VALID_MODES.includes(paramMode)) return paramMode;

    return (localStorage.getItem(VIEW_MODE_KEY) as ViewMode) || 'translation';
  });

  useEffect(() => {
    localStorage.setItem(VIEW_MODE_KEY, mode);
  }, [mode]);

  return [mode, setMode] as const;
}
