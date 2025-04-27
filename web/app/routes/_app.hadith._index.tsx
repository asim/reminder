export default function HadithIndex() {
  return (
    <div className='flex flex-col flex-1 p-8 max-w-3xl mx-auto'>
      <h1 className='text-4xl font-semibold mb-6 text-left'>Hadith</h1>

      <section className='mb-8'>
        <h2 className='text-xl font-medium mb-3'>What are Hadith?</h2>
        <p className='text-gray-700 mb-3'>
          Hadith are the recorded sayings, actions, and silent approvals of
          Prophet Muhammad (Peace Be Upon Him). These narrations were
          meticulously collected and authenticated by scholars in the early
          centuries of Islam, forming a vast body of literature that serves as
          the second most important source of Islamic teachings after the Quran.
        </p>
      </section>

      <section className='mb-8'>
        <h2 className='text-xl font-medium mb-3'>Hadith and the Quran</h2>
        <p className='text-gray-700 mb-3'>
          While the Quran is the direct word of God revealed to Prophet
          Muhammad, hadith complement the Quran by providing context,
          explanation, and practical examples. The Quran itself instructs
          Muslims to follow the Prophet's example, making hadith an integral
          part of understanding and practicing Islam.
        </p>
      </section>

      <section className='mb-8'>
        <h2 className='text-xl font-medium mb-3'>Collections of Hadith</h2>
        <p className='text-gray-700 mb-3'>
          Over time, scholars compiled various collections of hadith, carefully
          verifying their authenticity through a rigorous methodology. The most
          respected collections include Sahih Bukhari, Sahih Muslim, Sunan Abu
          Dawood, Jami at-Tirmidhi, Sunan an-Nasa'i, and Sunan Ibn Majah.
        </p>
        <p className='text-gray-700'>
          Hadith collection shown on this website are from Sahih Bukhari.
        </p>
      </section>

      <section className='mb-12 bg-gray-50 p-4 rounded-lg border border-gray-200'>
        <h2 className='text-lg font-medium mb-2 text-gray-800'>
          Navigating the Hadith
        </h2>
        <p className='text-gray-600'>
          Please select a book from the navigation menu on the left to begin
          exploring the hadith.
        </p>
      </section>
    </div>
  );
}
