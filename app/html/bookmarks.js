// Bookmarks management for Reminder app
// Stores bookmarks in localStorage organized by type (quran, hadith, names)

const STORAGE_KEY = 'reminder_bookmarks';

// Initialize bookmarks structure
function initBookmarks() {
  return {
    quran: {},
    hadith: {},
    names: {}
  };
}

// Get all bookmarks from localStorage
function getBookmarks() {
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

// Save bookmarks to localStorage
function saveBookmarks(bookmarks) {
  try {
    localStorage.setItem(STORAGE_KEY, JSON.stringify(bookmarks));
    return true;
  } catch (e) {
    console.error('Error saving bookmarks:', e);
    return false;
  }
}

// Add a bookmark
function addBookmark(type, key, label, url) {
  const bookmarks = getBookmarks();
  
  if (!bookmarks[type]) {
    bookmarks[type] = {};
  }
  
  bookmarks[type][key] = {
    label: label,
    url: url,
    timestamp: new Date().toISOString()
  };
  
  return saveBookmarks(bookmarks);
}

// Remove a bookmark
function deleteBookmark(type, key) {
  const bookmarks = getBookmarks();
  
  if (bookmarks[type] && bookmarks[type][key]) {
    delete bookmarks[type][key];
    return saveBookmarks(bookmarks);
  }
  
  return false;
}

// Check if a bookmark exists
function hasBookmark(type, key) {
  const bookmarks = getBookmarks();
  return bookmarks[type] && bookmarks[type][key] !== undefined;
}

// Toggle bookmark (add if not exists, remove if exists)
function toggleBookmark(type, key, label, url) {
  if (hasBookmark(type, key)) {
    deleteBookmark(type, key);
    return false; // removed
  } else {
    addBookmark(type, key, label, url);
    return true; // added
  }
}

// Update bookmark button state
function updateBookmarkButton(button, isBookmarked) {
  if (isBookmarked) {
    button.textContent = '★';
    button.title = 'Remove bookmark';
    button.classList.add('bookmarked');
  } else {
    button.textContent = '☆';
    button.title = 'Add bookmark';
    button.classList.remove('bookmarked');
  }
}

// Initialize bookmark buttons on a page
function initializeBookmarkButtons() {
  document.querySelectorAll('.bookmark-btn').forEach(button => {
    const type = button.dataset.type;
    const key = button.dataset.key;
    
    // Set initial state
    updateBookmarkButton(button, hasBookmark(type, key));
    
    // Add click handler
    button.addEventListener('click', function(e) {
      e.preventDefault();
      const label = this.dataset.label;
      const url = this.dataset.url;
      const isBookmarked = toggleBookmark(type, key, label, url);
      updateBookmarkButton(this, isBookmarked);
    });
  });
}

// Run initialization when DOM is ready
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initializeBookmarkButtons);
} else {
  initializeBookmarkButtons();
}
