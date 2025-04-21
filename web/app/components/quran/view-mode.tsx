import { cn } from '~/utils/classname';
import { Languages, Book, BookA } from 'lucide-react';

export type ViewMode = 'translation' | 'recitation' | 'english';

type ViewModeProps = {
  mode: ViewMode;
  onChange: (mode: ViewMode) => void;
};

const MODES: Array<{ id: ViewMode; icon: React.ReactNode }> = [
  { id: 'translation', icon: <Languages className='size-3.5' /> },
  { id: 'recitation', icon: <Book className='size-3.5' /> },
  { id: 'english', icon: <BookA className='size-3.5' /> },
];

export function ViewMode(props: ViewModeProps) {
  const mode = props.mode;
  const onChange = props.onChange;

  return (
    <div className='flex items-center justify-center gap-2 mb-8'>
      {MODES.map((viewMode) => (
        <button
          key={viewMode.id}
          type='button'
          onClick={() => onChange(viewMode.id)}
          className={cn(
            'flex items-center text-sm gap-1.5 px-3 py-1.5 rounded-lg cursor-pointer border border-black/10 hover:bg-black/5 capitalize',
            mode === viewMode.id && 'bg-black text-white hover:bg-black/90'
          )}
        >
          {viewMode.icon}
          {viewMode.id}
        </button>
      ))}
    </div>
  );
}
