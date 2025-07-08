export function meta() {
  return [
    { title: 'Names of Allah - Reminder' },
    {
      property: 'og:title',
      content: 'Names of Allah - Reminder',
    },
    {
      name: 'description',
      content:
        'Learn about the Names of Allah, the 99 attributes and qualities through which Muslims identify and connect with Allah',
    },
  ];
}

export default function NamesIndex() {
  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full max-w-4xl'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        Names of Allah
      </h1>

      <div className='space-y-4 sm:space-y-6'>
        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>
            What are the Names of Allah?
          </h2>
          <p className='text-sm sm:text-base text-gray-700'>
            The Names of Allah, known as Asma ul Husna (The Most Beautiful
            Names), are the 99 attributes and qualities through which Muslims
            identify and connect with Allah. These names are mentioned
            throughout the Quran and Hadith, and they describe Allah's perfect
            attributes, characteristics, and actions.
          </p>
        </section>

        <section>
          <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
            Key points about the Names of Allah:
          </p>
          <ul className='list-disc pl-4 sm:pl-6 space-y-1 sm:space-y-2 text-sm sm:text-base text-gray-700'>
            <li>
              Allah has 99 names mentioned in Islamic tradition, though His
              attributes are limitless.
            </li>
            <li>
              The Prophet Muhammad (Peace Be Upon Him) said: "Allah has
              ninety-nine names, one hundred minus one; whoever enumerates them
              will enter Paradise."
            </li>
            <li>
              Each name reveals a different aspect of Allah's nature and
              relationship with creation.
            </li>
            <li>
              Muslims often recite these names in devotional practices, seeking
              to understand the divine attributes.
            </li>
            <li>
              The greatest of these names is considered to be "Allah" itself,
              which encompasses all other attributes.
            </li>
          </ul>
        </section>

        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>
            Significance
          </h2>
          <p className='text-sm sm:text-base text-gray-700'>
            Understanding and contemplating the Names of Allah is central to
            developing a deeper connection with the Creator. These names help
            Muslims comprehend how Allah interacts with creation through His
            mercy, wisdom, power, and justice.
          </p>
          <p className='mt-1 sm:mt-2 text-sm sm:text-base text-gray-700'>
            The Quran encourages believers to call upon Allah through His
            beautiful names: "And to Allah belong the best names, so invoke Him
            by them" (Quran 7:180). Muslims often recite specific names in
            different situations – calling on Al-Shafi (The Healer) when sick,
            or Al-Razzaq (The Provider) when in need.
          </p>
        </section>

        <section className='mt-6 sm:mt-8 mb-8 sm:mb-12 bg-gray-50 p-3 sm:p-4 rounded-lg border border-gray-200'>
          <h2 className='text-base sm:text-lg font-medium mb-1 sm:mb-2 text-gray-800'>
            Exploring the Names of Allah
          </h2>
          <p className='text-sm sm:text-base text-gray-600'>
            Please select a name from the navigation menu to learn more about
            its meaning, significance, and presence in the Quran and Hadith.
            Understanding these names can deepen your spiritual connection and
            appreciation of Allah's attributes.
          </p>
        </section>
      </div>
    </div>
  );
}
