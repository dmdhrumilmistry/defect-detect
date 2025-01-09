import * as React from 'react';

import { cn } from '@/lib/utils';
import { borderClasses, focusVisibleClasses } from '@/styles/standard-classes';

const Input = React.forwardRef<HTMLInputElement, React.ComponentProps<'input'>>(
    ({ className, type, ...props }, ref) => {
        return (
            <input
                type={type}
                className={cn(
                    'flex h-9 w-full bg-transparent px-3 py-1 text-base shadow-sm transition-colors placeholder:text-slate-500 disabled:cursor-not-allowed disabled:opacity-50 md:text-sm',
                    'file:border-0 file:bg-transparent file:text-sm file:font-medium file:text-slate-950',
                    'dark:file:text-slate-50 dark:placeholder:text-slate-400',
                    borderClasses,
                    focusVisibleClasses,
                    className
                )}
                ref={ref}
                {...props}
            />
        );
    }
);
Input.displayName = 'Input';

export { Input };
