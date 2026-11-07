package app

import (
	"fmt"
	"github.com/asim/reminder/daily"
	"strings"
)

const Bookmarks = `
<div id="quran-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Quran</h2>
  <div id="quran-list" class="space-y-3"></div>
</div>

<div id="hadith-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Hadith</h2>
  <div id="hadith-list" class="space-y-3"></div>
</div>

<div id="names-bookmarks" class="mb-8">
  <h2 class="text-2xl font-semibold mb-4">Names</h2>
  <div id="names-list" class="space-y-3"></div>
</div>

<script>
  function loadBookmarks() {
    const bookmarks = JSON.parse(localStorage.getItem('reminder_bookmarks') || '{"quran":{},"hadith":{},"names":{}}');

    // Quran bookmarks
    const quranList = document.getElementById('quran-list');
    const quranKeys = Object.keys(bookmarks.quran || {});
    if (quranKeys.length > 0) {
      quranList.innerHTML = quranKeys.map(key => {
        const bookmark = bookmarks.quran[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('quran', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      quranList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }

    // Hadith bookmarks
    const hadithList = document.getElementById('hadith-list');
    const hadithKeys = Object.keys(bookmarks.hadith || {});
    if (hadithKeys.length > 0) {
      hadithList.innerHTML = hadithKeys.map(key => {
        const bookmark = bookmarks.hadith[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('hadith', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      hadithList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }

    // Names bookmarks
    const namesList = document.getElementById('names-list');
    const namesKeys = Object.keys(bookmarks.names || {});
    if (namesKeys.length > 0) {
      namesList.innerHTML = namesKeys.map(key => {
        const bookmark = bookmarks.names[key];
        return ` + "`" + `
        <div class="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors">
          <a href="${bookmark.url}" class="flex-grow text-blue-600 hover:text-blue-800">${bookmark.label}</a>
          <button onclick="removeBookmark('names', '${key}')" class="px-3 py-1 bg-red-500 text-white rounded hover:bg-red-600 transition-colors text-sm">Remove</button>
        </div>
      ` + "`" + `;
      }).join('');
    } else {
      namesList.innerHTML = '<p class="text-gray-500">No bookmarks yet</p>';
    }
  }

  function removeBookmark(type, key) {
    const bookmarks = JSON.parse(localStorage.getItem('reminder_bookmarks') || '{"quran":{},"hadith":{},"names":{}}');
    if (bookmarks[type] && bookmarks[type][key]) {
      delete bookmarks[type][key];
      localStorage.setItem('reminder_bookmarks', JSON.stringify(bookmarks));
      loadBookmarks();
    }
  }

  loadBookmarks();
</script>
`

const dailyTemplate = `
<h3 class="text-lg font-semibold mt-6 mb-2">Verse</h3>
<a href="%s" class="block p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors mb-4">%s</a>
<h3 class="text-lg font-semibold mt-6 mb-2">Hadith</h3>
<a href="%s" class="block p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors mb-4">%s</a>
<h3 class="text-lg font-semibold mt-6 mb-2">Name</h3>
<a href="%s" class="block p-4 bg-white border border-gray-200 rounded-lg hover:border-gray-400 transition-colors mb-4">%s</a>
<p class="text-sm text-gray-500 mt-4">Updated %s</p>
`

// HomeContent fills the existing lite home page with a reminder snapshot.
func HomeContent(r daily.Reminder) string {
	content := strings.ReplaceAll(Index, "{verse_link}", r.Links["verse"])
	content = strings.ReplaceAll(content, "{verse_text}", r.Verse)
	content = strings.ReplaceAll(content, "{hadith_link}", r.Links["hadith"])
	content = strings.ReplaceAll(content, "{hadith_text}", r.Hadith)
	content = strings.ReplaceAll(content, "{name_link}", r.Links["name"])
	return strings.ReplaceAll(content, "{name_text}", r.Name)
}

func DailyContent(r daily.Reminder) string {
	return fmt.Sprintf(dailyTemplate, r.Links["verse"], r.Verse, r.Links["hadith"], r.Hadith, r.Links["name"], r.Name, r.Updated)
}
