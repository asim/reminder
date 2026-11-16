export const tafsirEditions = [
  { id: 'jalalayn-en', label: 'Al-Jalalayn · English' },
  { id: 'original-en', label: 'Al-Mukhtasar · English' },
  { id: 'ibn-kathir-en', label: 'Ibn Kathir · English' },
] as const;
export type TafsirEdition = typeof tafsirEditions[number]['id'];
export const validTafsir = (value: unknown): value is TafsirEdition =>
  tafsirEditions.some(edition => edition.id === value);

