type ChapterHeaderProps = {
  title: string;
  translation: string;
  subtitle: string;
};

export function ChapterHeader(props: ChapterHeaderProps) {
  const { title, translation, subtitle } = props;

  return (
    <div className='text-center mb-12 mt-12'>
      <div className='text-xl text-gray-600'>{translation}</div>
      <h1 className='text-5xl font-bold my-4'>{title}</h1>
      <div className='text-xl text-gray-600'>{subtitle}</div>
    </div>
  );
}
