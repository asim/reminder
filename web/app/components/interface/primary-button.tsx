import { cn } from '~/utils/classname';
import { Slot } from '@radix-ui/react-slot';

type PrimaryButtonProps = React.ComponentProps<'button'> & {
  className?: string;
  disabled?: boolean;
  asChild?: boolean;
};

export function PrimaryButton(props: PrimaryButtonProps) {
  const { children, className, disabled, asChild, ...rest } = props;

  const Comp = asChild ? Slot : 'button';

  return (
    <Comp
      {...rest}
      disabled={disabled}
      className={cn(
        'bg-black flex flex-row items-center gap-2 hover:opacity-80 text-white px-4 py-2 rounded-lg',
        disabled && 'opacity-50 cursor-not-allowed pointer-events-none',
        className
      )}
    >
      {children}
    </Comp>
  );
}
