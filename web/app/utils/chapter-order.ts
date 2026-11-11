export type ChapterOrder = 'standard' | 'revelation';

// Traditional chronology published by Tanzil, based on the first revealed verses.
// Source: https://tanzil.net/docs/revelation_order (retrieved 2026-11-11).
// This is a chapter-list ordering, not a chronology of every verse.
export const revelationOrder = [
  96, 68, 73, 74, 1, 111, 81, 87, 92, 89, 93, 94,
  103, 100, 108, 102, 107, 109, 105, 113, 114, 112, 53, 80,
  97, 91, 85, 95, 106, 101, 75, 104, 77, 50, 90, 86,
  54, 38, 7, 72, 36, 25, 35, 19, 20, 56, 26, 27,
  28, 17, 10, 11, 12, 15, 6, 37, 31, 34, 39, 40,
  41, 42, 43, 44, 45, 46, 51, 88, 18, 16, 71, 14,
  21, 23, 32, 52, 67, 69, 70, 78, 79, 82, 84, 30,
  29, 83, 2, 8, 3, 33, 60, 4, 99, 57, 47, 13,
  55, 76, 65, 98, 59, 24, 22, 63, 58, 49, 66, 64,
  61, 62, 48, 5, 9, 110,
] as const;

const revelationRank = new Map<number, number>(revelationOrder.map((number, rank) => [number, rank]));

export function orderChapters<T extends { number: number }>(chapters: readonly T[], order: ChapterOrder): T[] {
  return [...chapters].sort((a, b) => order === 'revelation'
    ? (revelationRank.get(a.number) ?? Infinity) - (revelationRank.get(b.number) ?? Infinity)
    : a.number - b.number);
}
