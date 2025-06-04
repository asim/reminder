import React from 'react';

export function meta() {
  return [
    { title: 'Islam - Reminder' },
    {
      name: 'description',
      content: 'What is Islam? Learn about the faith, the Five Pillars of Islam, and the Articles of Faith.'
    }
  ];
}

export default function IslamPage() {
  return (
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-3xl overflow-y-auto px-5 py-5'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        Islam
      </h1>

      <section className='mb-4 sm:mb-8'>
        <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>What is Islam?</h2>
        <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
          Islam is a monotheistic faith centered on the belief in one God (Allah), who is merciful, compassionate, and just. Muslims believe that God sent prophets throughout history to guide humanity, with the final prophet being Muhammad (peace be upon him). Islam teaches submission to the will of God, striving for righteousness, and living a life of compassion, justice, and service to others. The Quran and the teachings of Prophet Muhammad (Hadith) form the foundation of Islamic belief and practice.
        </p>
      </section>

      <section className='mb-4 sm:mb-8'>
        <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>The Five Pillars of Islam</h2>
        <ol className='list-decimal list-inside space-y-2 text-base sm:text-lg'>
          <li><strong>Shahada (Faith):</strong> Declaring there is no god but Allah, and Muhammad is His Messenger.</li>
          <li><strong>Salah (Prayer):</strong> Performing the five daily prayers.</li>
          <li><strong>Zakat (Charity):</strong> Giving to those in need and supporting the community.</li>
          <li><strong>Sawm (Fasting):</strong> Fasting during the month of Ramadan.</li>
          <li><strong>Hajj (Pilgrimage):</strong> Making the pilgrimage to Mecca at least once if able.</li>
        </ol>
      </section>

      <section className='mb-8 sm:mb-12'>
        <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>The Six Articles of Faith</h2>
        <ul className='list-disc list-inside space-y-2 text-base sm:text-lg'>
          <li>Belief in Allah (God)</li>
          <li>Belief in the Angels</li>
          <li>Belief in the revealed Books</li>
          <li>Belief in the Prophets</li>
          <li>Belief in the Day of Judgment</li>
          <li>Belief in Divine Decree (Qadar)</li>
        </ul>
      </section>
    </div>
  );
}
