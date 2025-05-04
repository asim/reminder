type ChapterHeaderProps = {
  title: string;
  translation: string;
  subtitle: string;
};

export function ChapterHeader(props: ChapterHeaderProps) {
  const { title, translation, subtitle } = props;

  return (
    <div className='text-center mb-6 sm:mb-8 md:mb-12 mt-0 sm:mt-8 md:mt-12'>
      <div className='text-base sm:text-lg md:text-xl text-gray-600'>{translation}</div>
      <h1 className='text-3xl sm:text-4xl md:text-5xl font-bold my-2 sm:my-3 md:my-4'>{title}</h1>
      <div className='text-base sm:text-lg md:text-xl text-gray-600'>{subtitle}</div>
    </div>
  );
}
