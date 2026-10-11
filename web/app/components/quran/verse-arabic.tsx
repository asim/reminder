import { ClickableArabicWord } from './clickable-arabic-word';

type Verse = { number: number; arabic: string; words?: { arabic: string; english: string; transliteration: string }[] };

export function VerseArabic({ verse, wordByWord = false, transliteration = false }: {
  verse: Verse; wordByWord?: boolean; transliteration?: boolean;
}) {
  const number = String(verse.number).replace(/\d/g, d => String.fromCharCode(0x0660 + Number(d)));
  const annotated = (wordByWord || transliteration) && !!verse.words?.length;
  return (
    <div lang='ar' dir='rtl' className={`quran-arabic font-arabic text-3xl text-right leading-loose mb-3 sm:mb-4 ${annotated ? 'flex flex-wrap items-end' : ''}`}>
      {annotated ? verse.words!.map((word, i) => (
        <ClickableArabicWord key={i} arabic={word.arabic} transliteration={word.transliteration} className='inline-flex flex-col items-center mr-2 mb-2'>
          <span className='flex flex-row items-center'><span>{word.arabic}</span>{i === verse.words!.length - 1 && verse.number > 0 && <span className='mx-2' aria-label={`Verse ${verse.number}`}>{number}</span>}</span>
          {transliteration && <span lang='en' dir='ltr' className='font-sans text-xs mt-1 px-1 rounded bg-blue-50 text-blue-700'>{word.transliteration}</span>}
          {wordByWord && <span lang='en' dir='ltr' className='font-sans text-xs sm:text-sm mt-1 px-1 rounded bg-gray-100 text-gray-700'>{word.english}</span>}
        </ClickableArabicWord>
      )) : verse.arabic}
      {!annotated && verse.number > 0 && <span className='mx-2' aria-label={`Verse ${verse.number}`}>{number}</span>}
    </div>
  );
}
