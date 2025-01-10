import type { ChangeEvent } from 'react';
import { Search } from 'lucide-react';

import { cn } from '@/lib/utils';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { focusVisibleClasses, subtleIconStroke } from '@/styles/standard-classes';

type SearchBarProps = {
    searchQuery: string;
    className?: string;
    setSearchQuery: (value: string) => void;
};

export default function SearchBar({ searchQuery, className, setSearchQuery }: SearchBarProps) {
    return (
        <div className={cn('relative w-72', className)}>
            <Label htmlFor="search" className="sr-only">
                Search
            </Label>
            <Input
                id="search"
                type="text"
                autoComplete="off"
                value={searchQuery}
                onChange={(e: ChangeEvent<HTMLInputElement>) => setSearchQuery(e.target.value)}
                placeholder="Search for a project..."
                className={cn('bg-white w-full pl-8 text-sm shadow-none dark:bg-slate-950', focusVisibleClasses)}
            />
            <Search
                className={cn(
                    'pointer-events-none absolute left-2 top-1/2 size-4 -translate-y-1/2 select-none',
                    subtleIconStroke
                )}
            />
        </div>
    );
}
