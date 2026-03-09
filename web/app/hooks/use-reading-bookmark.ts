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
  quran?: ReadingBookmark;
  hadith?: ReadingBookmark;
  names?: ReadingBookmark;
}

function initReadingBookmarks(): ReadingBookmarksData {
  return {};
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
    return JSON.parse(stored);
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

  const setReadingBookmark = (
    type: ReadingBookmarkType,
    label: string,
    url: string,
    excerpt?: string
  ) => {
    const newBookmarks = { ...readingBookmarks };
    newBookmarks[type] = {
      label,
      url,
      timestamp: new Date().toISOString(),
      excerpt,
    };
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const clearReadingBookmark = (type: ReadingBookmarkType) => {
    const newBookmarks = { ...readingBookmarks };
    delete newBookmarks[type];
    setReadingBookmarks(newBookmarks);
    saveReadingBookmarks(newBookmarks);
  };

  const getReadingBookmark = (type: ReadingBookmarkType): ReadingBookmark | undefined => {
    return readingBookmarks[type];
  };

  const isReadingBookmark = (type: ReadingBookmarkType, url: string): boolean => {
    return readingBookmarks[type]?.url === url;
  };

  return {
    readingBookmarks,
    setReadingBookmark,
    clearReadingBookmark,
    getReadingBookmark,
    isReadingBookmark,
  };
}
