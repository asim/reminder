import { useState } from 'react';

const STORAGE_KEY = 'reminder_reading_bookmarks';

export type ReadingBookmarkType = 'quran' | 'hadith' | 'names';

export interface ReadingBookmark {
  label: string;
  url: string;
  timestamp: string;
  excerpt?: string;
}

export interface ReadingBookmarksData {
  quran: ReadingBookmark[];
  hadith: ReadingBookmark[];
  names: ReadingBookmark[];
}

function initReadingBookmarks(): ReadingBookmarksData {
  return { quran: [], hadith: [], names: [] };
}

/**
 * Migrate from old single-bookmark format:
 *   { quran?: ReadingBookmark, hadith?: ReadingBookmark, names?: ReadingBookmark }
 * to new multi-bookmark format:
 *   { quran: ReadingBookmark[], hadith: ReadingBookmark[], names: ReadingBookmark[] }
 */
function migrateIfNeeded(raw: Record<string, unknown>): ReadingBookmarksData {
  const result = initReadingBookmarks();

  for (const type of ['quran', 'hadith', 'names'] as ReadingBookmarkType[]) {
    const value = raw[type];
    if (!value) continue;

    if (Array.isArray(value)) {
      // Already in new format
      result[type] = value as ReadingBookmark[];
    } else if (typeof value === 'object' && 'url' in (value as Record<string, unknown>)) {
      // Old single-bookmark format — wrap in array
      result[type] = [value as ReadingBookmark];
    }
  }

  return result;
}

function getStoredReadingBookmarks(): ReadingBookmarksData {
  if (typeof window === 'undefined') {
    return initReadingBookmarks();
  }

  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) {
      return initReadingBookmarks();
    }
    const parsed = JSON.parse(stored);
    return migrateIfNeeded(parsed);
  } catch (e) {
    console.error('Error reading reading bookmarks:', e);
    return initReadingBookmarks();
  }
}

function saveReadingBookmarks(bookmarks: ReadingBookmarksData): boolean {
  if (typeof window === 'undefined') {
    return false;
  }

  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(bookmarks));
    return true;
  } catch (e) {
    console.error('Error saving reading bookmarks:', e);
    return false;
  }
}

export function useReadingBookmark() {
  const [readingBookmarks, setReadingBookmarks] = useState<ReadingBookmarksData>(getStoredReadingBookmarks);

  const addReadingBookmark = (
    type: ReadingBookmarkType,
    label: string,
    url: string,
    excerpt?: string
  ) => {
    const newBookmarks = { ...readingBookmarks };
    // Remove any existing bookmark with the same URL to avoid duplicates
    newBookmarks[type] = newBookmarks[type].filter(b => b.url !== url);
    // Add new bookmark at the front (most recent first)
    newBookmarks[type] = [
      { label, url, timestamp: new Date().toISOString(), excerpt },
      ...newBookmarks[type],
    ];
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const removeReadingBookmark = (type: ReadingBookmarkType, url: string) => {
    const newBookmarks = { ...readingBookmarks };
    newBookmarks[type] = newBookmarks[type].filter(b => b.url !== url);
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const clearReadingBookmark = (type: ReadingBookmarkType) => {
    const newBookmarks = { ...readingBookmarks };
    newBookmarks[type] = [];
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const isReadingBookmark = (type: ReadingBookmarkType, url: string): boolean => {
    return readingBookmarks[type].some(b => b.url === url);
  };

  // Backward-compatible: return the most recent bookmark for a type
  const getReadingBookmark = (type: ReadingBookmarkType): ReadingBookmark | undefined => {
    return readingBookmarks[type][0];
  };

  // Legacy single-bookmark setter for backward compat (wraps addReadingBookmark)
  const setReadingBookmark = addReadingBookmark;

  return {
    readingBookmarks,
    setReadingBookmark,
    addReadingBookmark,
    removeReadingBookmark,
    clearReadingBookmark,
    getReadingBookmark,
    isReadingBookmark,
  };
}
