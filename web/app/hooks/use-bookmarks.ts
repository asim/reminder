import { useEffect, useState } from 'react';

const STORAGE_KEY = 'reminder_bookmarks';

export type BookmarkType = 'quran' | 'hadith' | 'names';

export interface Bookmark {
  label: string;
  url: string;
  timestamp: string;
  excerpt?: string;
}

export interface BookmarksData {
  quran: Record<string, Bookmark>;
  hadith: Record<string, Bookmark>;
  names: Record<string, Bookmark>;
}

function initBookmarks(): BookmarksData {
  return {
    quran: {},
    hadith: {},
    names: {},
  };
}

function getStoredBookmarks(): BookmarksData {
  if (typeof window === 'undefined') {
    return initBookmarks();
  }

  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) {
      return initBookmarks();
    }
    return JSON.parse(stored);
  } catch (e) {
    console.error('Error reading bookmarks:', e);
    return initBookmarks();
  }
}

function saveBookmarks(bookmarks: BookmarksData): boolean {
  if (typeof window === 'undefined') {
    return false;
  }

  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(bookmarks));
    return true;
  } catch (e) {
    console.error('Error saving bookmarks:', e);
    return false;
  }
}

export function useBookmarks() {
  const [bookmarks, setBookmarks] = useState<BookmarksData>(getStoredBookmarks);

  const addBookmark = (
    type: BookmarkType,
    key: string,
    label: string,
    url: string,
    excerpt?: string
  ) => {
    const newBookmarks = { ...bookmarks };
    newBookmarks[type][key] = {
      label,
      url,
      timestamp: new Date().toISOString(),
      excerpt,
    };
    setBookmarks(newBookmarks);
    saveBookmarks(newBookmarks);
  };

  const removeBookmark = (type: BookmarkType, key: string) => {
    const newBookmarks = { ...bookmarks };
    delete newBookmarks[type][key];
    setBookmarks(newBookmarks);
    saveBookmarks(newBookmarks);
  };

  const hasBookmark = (type: BookmarkType, key: string): boolean => {
    return bookmarks[type][key] !== undefined;
  };

  const toggleBookmark = (
    type: BookmarkType,
    key: string,
    label: string,
    url: string,
    excerpt?: string
  ): boolean => {
    if (hasBookmark(type, key)) {
      removeBookmark(type, key);
      return false;
    } else {
      addBookmark(type, key, label, url, excerpt);
      return true;
    }
  };

  const getAllBookmarks = (): BookmarksData => {
    return bookmarks;
  };

  return {
    bookmarks,
    addBookmark,
    removeBookmark,
    hasBookmark,
    toggleBookmark,
    getAllBookmarks,
  };
}
