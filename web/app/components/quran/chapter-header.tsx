type ChapterHeaderProps = {
  title: string;
  subtitle: string;
};

export function ChapterHeader(props: ChapterHeaderProps) {
  const { title, subtitle } = props;

  return (
    <div className='text-center mb-12 mt-12'>
      <h1 className='text-4xl font-bold mb-2'>{title}</h1>
      <div className='text-xl text-gray-600'>{subtitle}</div>
    </div>
  );
}
