export function readingMeta(kind: string, { data, params }: any) {
  const record = data?.readingQueries?.[0]?.data;
  let title = 'Reminder';
  let description = 'Read the Quran, Hadith and Names of Allah.';
  if (kind === 'chapter' || kind === 'verse') {
    const ref = `${params.chapterNumber}${kind === 'verse' ? `:${params.verseNumber}` : ''}`;
    title = `${record?.name || 'Quran'} ${ref} - Quran - Reminder`;
    description = kind === 'verse' ? record?.verses?.[0]?.text : `Read ${record?.name || 'this chapter'} (${record?.english || ref}) in Arabic and English with tafsir.`;
  } else if (kind === 'name') {
    title = `${record?.english || 'Names of Allah'} — ${record?.meaning || params.nameNumber} - Reminder`;
    description = record?.summary || record?.description;
  } else {
    const collection = data?.readingQueries?.[0]?.queryKey?.[1] === 'muslim' ? 'Sahih Muslim' : 'Sahih Bukhari';
    title = `${record?.name || `Book ${params.book}`} - ${collection} - Reminder`;
    description = `Read ${record?.name || 'this book'} from ${collection} in Arabic and English.`;
  }
  return [{ title }, { name: 'description', content: description || title }, { property: 'og:title', content: title }];
}
