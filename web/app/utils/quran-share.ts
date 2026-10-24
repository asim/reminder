import type { ViewMode } from '~/components/quran/view-mode';

export interface QuranViewSettings {
  mode: ViewMode;
  translation?: 'khattab' | 'asad' | 'sahih';
  wordByWord: boolean;
  commentary: boolean;
  transliteration?: boolean;
}

export function buildQuranShareUrl(
  basePath: string,
  settings: QuranViewSettings
): string {
  const params = new URLSearchParams();

  // Explicit false values keep shared links independent of saved preferences.
  params.set('mode', settings.mode);
  params.set('translation', settings.translation || 'khattab');
  params.set('wbw', settings.wordByWord ? '1' : '0');
  params.set('commentary', settings.commentary ? '1' : '0');
  params.set('transliteration', settings.transliteration ? '1' : '0');

  const qs = params.toString();
  return qs ? `${basePath}?${qs}` : basePath;
}
