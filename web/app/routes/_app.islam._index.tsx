import React from 'react';

export function meta() {
  return [
    { title: 'Islam - The Pillars and Core Concepts' },
    {
      name: 'description',
      content: 'Learn about the Five Pillars of Islam and the core concepts of the faith in a simple, accessible way.'
    }
  ];
}

export default function IslamPage() {
  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        The Pillars of Islam
      </h1>

      <section className="mb-4 sm:mb-8">
        <h2 className="text-lg sm:text-xl font-medium mb-2 sm:mb-3">
          The Five Pillars of Islam
        </h2>
        <ol className="list-decimal list-inside space-y-2 text-base sm:text-lg">
          <li><strong>Shahada (Faith):</strong> Declaring there is no god but Allah, and Muhammad is His Messenger.</li>
          <li><strong>Salah (Prayer):</strong> Performing the five daily prayers.</li>
          <li><strong>Zakat (Charity):</strong> Giving to those in need and supporting the community.</li>
          <li><strong>Sawm (Fasting):</strong> Fasting during the month of Ramadan.</li>
          <li><strong>Hajj (Pilgrimage):</strong> Making the pilgrimage to Mecca at least once if able.</li>
        </ol>
      </section>

      <section className="mb-4 sm:mb-8">
        <h2 className="text-lg sm:text-xl font-medium mb-2 sm:mb-3">
          Core Concepts
        </h2>
        <ul className="list-disc list-inside space-y-2 text-base sm:text-lg">
          <li><strong>Tawhid:</strong> The oneness of God. God is unique, without partner, child, or equal.</li>
          <li><strong>Prophethood:</strong> God sent messengers to guide humanity.</li>
          <li><strong>Revelation:</strong> God’s guidance is delivered through scriptures like the Quran.</li>
          <li><strong>Worship:</strong> Acts of devotion are directed to God alone.</li>
          <li><strong>Accountability:</strong> Every person is responsible for their actions.</li>
        </ul>
      </section>

      <section className="mb-8 sm:mb-12 bg-gray-50 p-3 sm:p-4 rounded-lg border border-gray-200">
        <h2 className="text-base sm:text-lg font-medium mb-1 sm:mb-2 text-gray-800">
          Why These Pillars Matter
        </h2>
        <p className="text-sm sm:text-base text-gray-700">
          Understanding these pillars and concepts helps build a strong foundation for faith and provides a welcoming entry point for anyone interested in learning more about Islam.
        </p>
      </section>
    </div>
  );
}
