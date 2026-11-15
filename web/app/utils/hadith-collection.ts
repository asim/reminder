export type HadithCollection = 'bukhari' | 'muslim';
export const collectionFromPath = (path: string): HadithCollection => /^\/hadith\/muslim(?:\/|$)/.test(path) ? 'muslim' : 'bukhari';
export const hadithPath = (collection: HadithCollection) => collection === 'muslim' ? '/hadith/muslim' : '/hadith';
export const collectionTitle = (collection: HadithCollection) => collection === 'muslim' ? 'Sahih Muslim' : 'Sahih Bukhari';
