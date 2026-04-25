import type { ViewMode } from '~/components/quran/view-mode';

export interface QuranViewSettings {
  mode: ViewMode;
  wordByWord: boolean;
  commentary: boolean;
}

export function buildQuranShareUrl(
  basePath: string,
  settings: QuranViewSettings
): string {
  const params = new URLSearchParams();

  if (settings.mode !== 'translation') {
    params.set('mode', settings.mode);
  }
  if (settings.wordByWord) {
    params.set('wbw', '1');
  }
  if (settings.commentary) {
    params.set('commentary', '1');
  }

  const qs = params.toString();
  return qs ? `${basePath}?${qs}` : basePath;
}
