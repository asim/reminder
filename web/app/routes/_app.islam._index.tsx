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
    <div className='flex flex-col flex-1 p-0 lg:p-8 mx-auto w-full lg:max-w-4xl overflow-y-auto px-5 py-5'>
      <h1 className='text-2xl sm:text-3xl md:text-4xl font-semibold mb-4 sm:mb-6 text-left'>
        Islam
      </h1>

      <section className='mb-4 sm:mb-8'>
        <h2 className='text-lg sm:text-xl font-medium mb-2 sm:mb-3'>What is Islam?</h2>
        <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
          Islam is a monotheistic faith centered on the belief in one God (Allah), who is merciful, compassionate, and just. Muslims believe that God sent prophets throughout history to guide humanity, with the final prophet being Muhammad (peace be upon him). Islam teaches submission to the will of God, striving for righteousness, and living a life of compassion, justice, and service to others. The Quran and the teachings of Prophet Muhammad (Hadith) form the foundation of Islamic belief and practice.
        </p>
      </section>

      <div className='space-y-4 sm:space-y-6'>
        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>Pillars of Islam</h2>
          <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
            The Pillars of Islam are five essential acts of worship and practice that form the foundation of a Muslim's faith and actions. They are considered the core duties that every Muslim is expected to observe, serving as a guide for living a life devoted to God and in service to others.
          </p>
          <ol className='list-decimal pl-4 sm:pl-6 space-y-1 sm:space-y-2 text-sm sm:text-base text-gray-700'>
            <li><strong>Shahada (Faith):</strong> The declaration of faith. A Muslim testifies that there is no god but Allah, and Muhammad is His Messenger. This simple statement affirms the core belief in the oneness of God and the acceptance of Muhammad as His final prophet. It is the entry point into Islam and shapes a Muslim's worldview and purpose in life.</li>
            <li><strong>Salah (Prayer):</strong> Muslims are required to pray five times a day at prescribed times. These prayers are a direct link between the worshipper and God, providing structure to the day and serving as a reminder of the importance of faith, gratitude, and discipline in daily life.</li>
            <li><strong>Zakat (Charity):</strong> Muslims who are able must give a portion of their wealth each year to those in need. This practice purifies wealth, fosters social responsibility, and helps reduce inequality by supporting the less fortunate in the community.</li>
            <li><strong>Sawm (Fasting):</strong> During the month of Ramadan, Muslims fast from dawn to sunset. Fasting teaches self-control, empathy for the hungry, and spiritual reflection. It is a time for increased worship, charity, and community.</li>
            <li><strong>Hajj (Pilgrimage):</strong> Muslims who are physically and financially able must perform the pilgrimage to Mecca at least once in their lifetime. Hajj is a profound spiritual journey that unites Muslims from around the world in worship and equality before God.</li>
          </ol>
        </section>

        <section>
          <h2 className='text-lg sm:text-xl font-medium mb-1 sm:mb-2'>Articles of Faith</h2>
          <p className='text-sm sm:text-base text-gray-700 mb-2 sm:mb-4'>
            The Articles of Faith are six core beliefs that every Muslim holds. These beliefs define a Muslim's understanding of the world, the unseen, and the relationship between God, humanity, and the universe. They provide the spiritual and theological framework for Islamic faith.
          </p>
          <ul className='list-disc pl-4 sm:pl-6 space-y-1 sm:space-y-2 text-sm sm:text-base text-gray-700'>
            <li><strong>Belief in Allah (God):</strong> Muslims believe in one, unique, all-powerful, and merciful God who created and sustains everything. God has no partners, children, or equals.</li>
            <li><strong>Belief in the Angels:</strong> Angels are spiritual beings created by God who carry out His commands, record human deeds, and deliver revelations to the prophets.</li>
            <li><strong>Belief in the revealed Books:</strong> Muslims believe that God sent guidance to humanity through scriptures, including the Torah, Psalms, Gospel, and finally the Quran, which is considered the final and complete revelation.</li>
            <li><strong>Belief in the Prophets:</strong> God sent many prophets to guide people, including Adam, Noah, Abraham, Moses, Jesus, and Muhammad (peace be upon them all). Muhammad is regarded as the last prophet.</li>
            <li><strong>Belief in the Day of Judgment:</strong> Muslims believe that everyone will be resurrected for judgment by God. Each person will be held accountable for their actions and rewarded or punished accordingly.</li>
            <li><strong>Belief in Divine Decree (Qadar):</strong> Muslims believe that God has knowledge and control over everything that happens. While humans have free will, everything ultimately occurs by God's will and wisdom.</li>
          </ul>
        </section>
      </div>
    </div>
  );
}
