export function meta() {
  return [
    { title: 'Quran - Reminder' },
    {
      property: 'og:title',
      content: 'Quran - Reminder',
    },
    {
      name: 'description',
      content: 'Read the Quran, the holy book of Islam',
    },
  ];
}

export default function QuranIndex() {
  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        The Quran
      </h1>

      <div className='space-y-4 sm:space-y-6'>
        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>
            What is the Quran?
          </h2>
          <p className='text-sm sm:text-base text-gray-700'>
            The Quran is the holy book of Islam, believed to be the word of
            Allah as revealed to the Prophet Muhammad through the angel Gabriel.
            It serves as a guidance and a message for all of humanity,
            containing verses that are perfectly explained and revealed in
            Arabic. It is considered a complete and clear book, serving as a
            guide, a reminder, and a mercy for those who have faith.
          </p>
        </section>

        <section>
          <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
            Some key facts about the Quran:
          </p>
          <ul className='list-disc pl-4 sm:pl-6 space-y-1 sm:space-y-2 text-sm sm:text-base text-gray-700'>
            <li>
              Revealed to Prophet Muhammad (Peace Be Upon Him) over a period of
              approximately 23 years, from 610 CE to 632 CE.
            </li>
            <li>Consists of 114 chapters (surahs)</li>
            <li>Contains approximately 6,236 verses (ayat).</li>
            <li>
              Written in classical Arabic, though translations are available in
              virtually every language.
            </li>
            <li>Memorized in its entirety by millions of Muslims worldwide.</li>
            <li>
              The Quran has been preserved in its original form since its
              revelation over 1400 years ago, without any additions or deletions
              to its text.
            </li>
          </ul>
        </section>

        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>
            Structure
          </h2>
          <p className='text-sm sm:text-base text-gray-700'>
            Each chapter (surah) of the Quran has a name, often derived from a
            notable word or theme within it. Surahs are classified as either
            Meccan (revealed in Mecca) or Medinan (revealed after migration to
            Medina), with each having distinct themes and characteristics.
          </p>
          <p className='mt-1 sm:mt-2 text-sm sm:text-base text-gray-700'>
            The first surah, Al-Fatihah (The Opening), is a short prayer that
            Muslims recite in their daily prayers. The longest surah is the
            second one, Al-Baqarah (The Cow), with 286 verses.
          </p>
        </section>

        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>
            Significance
          </h2>
          <p className='text-sm sm:text-base text-gray-700'>
            For Muslims, the Quran is more than just a religious text—it's a
            comprehensive guide to life. It addresses spiritual, social, legal,
            and personal matters, offering guidance on how to live according to
            divine will.
          </p>
          <p className='mt-1 sm:mt-2 text-sm sm:text-base text-gray-700'>
            Unlike many religious texts, the Quran is meant to be recited aloud,
            and its melodious nature is considered part of its miracle. The art
            of Quranic recitation (tajweed) is highly developed in Islamic
            tradition.
          </p>
        </section>

        <section className='mt-6 sm:mt-8 mb-8 sm:mb-12 bg-gray-50 p-3 sm:p-4 rounded-lg border border-gray-200'>
          <h2 className='text-base sm:text-lg font-medium mb-1 sm:mb-2 text-gray-800'>
            Navigating the Quran
          </h2>
          <p className='text-sm sm:text-base text-gray-600'>
            Please select a surah from the navigation menu to begin exploring
            the Quran. Each surah contains verses (ayat) that you can read in
            Arabic as well as translations.
          </p>
        </section>
      </div>
    </div>
  );
}
