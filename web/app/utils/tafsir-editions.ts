export const tafsirEditions = [
  { id: 'jalalayn-en', label: 'Al-Jalalayn · English' },
  { id: 'original-en', label: 'Original commentary · English (unattributed)' },
] as const;
export type TafsirEdition = typeof tafsirEditions[number]['id'];
export const validTafsir = (value: unknown): value is TafsirEdition =>
  tafsirEditions.some(edition => edition.id === value);

