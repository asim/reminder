import type { ReactNode } from 'react';
type ChapterHeaderProps = {
  title: string;
  translation: string;
  subtitle: string;
  actions?: ReactNode;
};

export function ChapterHeader(props: ChapterHeaderProps) {
  const { title, translation, subtitle, actions } = props;

  return (
    <div className='text-center mb-6 sm:mb-8 md:mb-12 mt-0 sm:mt-8 md:mt-12'>
      <div className='text-base sm:text-lg md:text-xl text-gray-600'>{translation}</div>
      <h1 className='text-3xl sm:text-4xl md:text-5xl font-bold my-2 sm:my-3 md:my-4'>{title}</h1>
      <div className='grid grid-cols-[1fr_auto_1fr] items-center gap-3'>
        <span aria-hidden='true' />
        <span className='text-base sm:text-lg md:text-xl text-gray-600'>{subtitle}</span>
        {actions && <div className='flex items-center gap-2 justify-self-start'>{actions}</div>}
      </div>
    </div>
  );
}
