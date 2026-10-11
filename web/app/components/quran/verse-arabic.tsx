import { ClickableArabicWord } from './clickable-arabic-word';

type Verse = { number: number; arabic: string; words?: { arabic: string; english: string; transliteration: string }[] };

export function VerseArabic({ verse, wordByWord = false, transliteration = false }: {
  verse: Verse; wordByWord?: boolean; transliteration?: boolean;
}) {
  const number = String(verse.number).replace(/\d/g, d => String.fromCharCode(0x0660 + Number(d)));
  const annotated = (wordByWord || transliteration) && !!verse.words?.length;
  return (
    <div lang='ar' dir='rtl' className={`quran-arabic font-arabic text-3xl text-right leading-loose mb-4 ${annotated ? 'flex flex-wrap items-end gap-x-3 gap-y-2' : ''}`}>
      {annotated ? verse.words!.map((word, i) => (
        <ClickableArabicWord key={i} arabic={word.arabic} transliteration={word.transliteration} className='inline-flex flex-col items-center'>
          <span>{word.arabic}</span>
          {transliteration && <span lang='en' dir='ltr' className='font-sans text-sm text-blue-700'>{word.transliteration}</span>}
          {wordByWord && <span lang='en' dir='ltr' className='font-sans text-sm text-gray-600'>{word.english}</span>}
        </ClickableArabicWord>
      )) : verse.arabic}
      {verse.number > 0 && <span className='mx-2' aria-label={`Verse ${verse.number}`}>﴿{number}﴾</span>}
    </div>
  );
}
