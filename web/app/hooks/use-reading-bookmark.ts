import { useEffect, useState } from 'react';

const STORAGE_KEY = 'reminder_reading_bookmarks';
const CHANGE_EVENT = 'reading-bookmarks-changed';

// Display preferences do not change the saved reading position.
export function readingPosition(url: string): string {
  try {
    const parsed = new URL(url, 'https://reminder.dev');
    const path = parsed.pathname.replace(/\/$/, '');
    const verse = path.match(/^\/quran\/(\d+)\/(\d+)$/);
    if (verse) return `/quran/${verse[1]}#${verse[2]}`;
    return path + parsed.hash;
  } catch {
    return url;
  }
}


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
    window.dispatchEvent(new Event(CHANGE_EVENT));
    return true;
  } catch (e) {
    console.error('Error saving reading bookmarks:', e);
    return false;
  }
}

export function useReadingBookmark() {
  const [readingBookmarks, setReadingBookmarks] = useState<ReadingBookmarksData>(getStoredReadingBookmarks);

  useEffect(() => {
    const refresh = () => setReadingBookmarks(getStoredReadingBookmarks());
    const onStorage = (event: StorageEvent) => {
      if (event.key === STORAGE_KEY || event.key === null) refresh();
    };
    window.addEventListener(CHANGE_EVENT, refresh);
    window.addEventListener('storage', onStorage);
    refresh();
    return () => {
      window.removeEventListener(CHANGE_EVENT, refresh);
      window.removeEventListener('storage', onStorage);
    };
  }, []);

  const addReadingBookmark = (
    type: ReadingBookmarkType,
    label: string,
    url: string,
    excerpt?: string
  ) => {
    const newBookmarks = getStoredReadingBookmarks();
    // Remove any existing bookmark with the same URL to avoid duplicates
    newBookmarks[type] = newBookmarks[type].filter(b => readingPosition(b.url) !== readingPosition(url));
    // Add new bookmark at the front (most recent first)
    newBookmarks[type] = [
      { label, url, timestamp: new Date().toISOString(), excerpt },
      ...newBookmarks[type],
    ];
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const removeReadingBookmark = (type: ReadingBookmarkType, url: string) => {
    const newBookmarks = getStoredReadingBookmarks();
    newBookmarks[type] = newBookmarks[type].filter(b => readingPosition(b.url) !== readingPosition(url));
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const clearReadingBookmark = (type: ReadingBookmarkType) => {
    const newBookmarks = getStoredReadingBookmarks();
    newBookmarks[type] = [];
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const isReadingBookmark = (type: ReadingBookmarkType, url: string): boolean => {
    return readingBookmarks[type].some(b => readingPosition(b.url) === readingPosition(url));
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
