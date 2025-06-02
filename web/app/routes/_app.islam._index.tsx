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
    <div className="flex flex-col items-center justify-center min-h-screen px-4 sm:px-6 md:px-8 py-12 sm:py-16 bg-white overflow-auto">
      <h1 className="text-4xl sm:text-5xl md:text-6xl font-bold mb-4 text-center">The Pillars of Islam</h1>
      <p className="text-gray-600 text-lg sm:text-xl mb-8 text-center max-w-2xl">
        Welcome! This page introduces the foundational pillars and core concepts of Islam in a simple, accessible way for everyone.
      </p>
      <section className="mb-8 w-full max-w-2xl">
        <h2 className="text-2xl font-semibold mb-4">The Five Pillars of Islam</h2>
        <ol className="list-decimal list-inside space-y-2 text-lg">
          <li><strong>Shahada (Faith):</strong> Declaring there is no god but Allah, and Muhammad is His Messenger.</li>
          <li><strong>Salah (Prayer):</strong> Performing the five daily prayers.</li>
          <li><strong>Zakat (Charity):</strong> Giving to those in need and supporting the community.</li>
          <li><strong>Sawm (Fasting):</strong> Fasting during the month of Ramadan.</li>
          <li><strong>Hajj (Pilgrimage):</strong> Making the pilgrimage to Mecca at least once if able.</li>
        </ol>
      </section>
      <section className="mb-8 w-full max-w-2xl">
        <h2 className="text-2xl font-semibold mb-4">Core Concepts</h2>
        <ul className="list-disc list-inside space-y-2 text-lg">
          <li><strong>Tawhid:</strong> The oneness of God. God is unique, without partner, child, or equal.</li>
          <li><strong>Prophethood:</strong> God sent messengers to guide humanity.</li>
          <li><strong>Revelation:</strong> God’s guidance is delivered through scriptures like the Quran.</li>
          <li><strong>Worship:</strong> Acts of devotion are directed to God alone.</li>
          <li><strong>Accountability:</strong> Every person is responsible for their actions.</li>
        </ul>
      </section>
      <section className="w-full max-w-2xl">
        <h2 className="text-2xl font-semibold mb-4">Why These Pillars Matter</h2>
        <p className="text-lg">
          Understanding these pillars and concepts helps build a strong foundation for faith and provides a welcoming entry point for anyone interested in learning more about Islam.
        </p>
      </section>
    </div>
  );
}
