import type { TProject } from '@/types';
import { Link } from 'react-router-dom';
import { formatDistanceToNow } from 'date-fns';
import { Ellipsis, Trash } from 'lucide-react';

// shadcn/ui components
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import {
    DropdownMenu,
    DropdownMenuTrigger,
    DropdownMenuContent,
    DropdownMenuItem,
} from '@/components/ui/dropdown-menu';

// utilities
import { cn } from '@/lib/utils';
import { borderClasses, focusVisibleClasses, subtleIconStroke, subtleText } from '@/styles/standard-classes';

type ProjectCardProps = {
    project: TProject;
    setProjectMarkedForDeletion: (value: TProject) => void;
};

export default function ProjectCard({ project, setProjectMarkedForDeletion }: ProjectCardProps) {
    return (
        <div
            className={cn(
                'relative min-w-48 flex flex-col transition-all shadow-base hover:shadow-md',
                borderClasses,
                focusVisibleClasses
            )}
        >
            <Link
                to={`/projects/${project.id}`}
                className="absolute inset-0 z-10 cursor-pointer overflow-hidden rounded-md"
            >
                <span className="sr-only">View Project</span>
            </Link>

            <h3 className="px-3 py-6 text-sm font-medium text-slate-900">{project.title}</h3>
            <Separator className="mx-3 w-auto"></Separator>
            <div className="flex justify-between h-11 items-center px-3">
                <span className={subtleText}>
                    Updated {formatDistanceToNow(project?.meta?.updatedAt, { addSuffix: true })}
                </span>

                {/* Project actions */}
                <DropdownMenu>
                    <DropdownMenuTrigger asChild className="z-10">
                        <Button size="icon" variant="ghost" className="h-7 w-7 [&>svg]:hover:stroke-slate-900">
                            <Ellipsis className={subtleIconStroke} />
                            <span className="sr-only">Open project actions</span>
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end" className="p-1.5 z-50 max-w-[300px]">
                        <DropdownMenuItem className="p-0 block">
                            <Button
                                variant="ghost"
                                className="px-2 w-full justify-start h-8 text-red-700 hover:bg-red-50 hover:text-red-700"
                                onClick={() => setProjectMarkedForDeletion(project)}
                            >
                                <Trash />
                                Delete Project
                                <span className="sr-only">Delete project action</span>
                            </Button>
                        </DropdownMenuItem>
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
}
