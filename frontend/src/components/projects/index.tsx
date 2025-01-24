import type { ProjectsDataLoader, TProject } from '@/types';
import { useEffect, useState } from 'react';
import { useFetcher, useLoaderData } from 'react-router-dom';
import { Loader2 } from 'lucide-react';

// shadcn/ui components
import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogFooter,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogCancel,
} from '@/components/ui/alert-dialog';
import { Button } from '@/components/ui/button';

// internal components
import CreateProject from './create-project';
import SearchBar from '@/components/shared/search-bar';
import ProjectCard from '@/components/shared/project-card';

export default function Projects() {
    const { projects } = useLoaderData() as ProjectsDataLoader;
    const fetcher = useFetcher();
    const [searchQuery, setSearchQuery] = useState('');
    const [projectMarkedForDeletion, setProjectMarkedForDeletion] = useState<TProject | null>(null);
    console.info('[COMP] Projects :: ', projects, fetcher);

    /**
     * :: Perf Note ::
     * Below filter will execute on each re-render hence this will become a bottleneck if we have large number of projects.
     */
    const filteredProjects = projects.filter((project) =>
        project.title.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const clearProjectMarkedForDeletion = () => {
        if (fetcher.state === 'idle' && Boolean(projectMarkedForDeletion)) setProjectMarkedForDeletion(null);
    };
    const handleProjectDeletion = () => {
        if (!projectMarkedForDeletion) return;
        fetcher.submit(null, {
            method: 'delete',
            action: `/projects/${projectMarkedForDeletion.id}`,
            encType: 'application/json',
        });
    };

    // if fetcher changes state then clear the project state is holding
    useEffect(() => {
        clearProjectMarkedForDeletion();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [fetcher.state]);

    return (
        <>
            <div className="min-h-screen mx-auto max-w-[1392px] w-full p-4">
                <div className="flex items-center justify-end gap-4">
                    <SearchBar searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
                    <CreateProject />
                </div>
                <div className="w-full pt-4 grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
                    {filteredProjects.map((project) => (
                        <ProjectCard
                            key={project.id}
                            project={project}
                            setProjectMarkedForDeletion={setProjectMarkedForDeletion}
                        />
                    ))}
                </div>
            </div>

            {/* confirm project deletion alert popup */}
            <AlertDialog open={Boolean(projectMarkedForDeletion)} onOpenChange={clearProjectMarkedForDeletion}>
                <AlertDialogContent className="max-w-[80%] sm:max-w-lg">
                    <AlertDialogHeader>
                        <AlertDialogTitle>Delete Project</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you sure you want to delete &quot;{projectMarkedForDeletion?.title}&quot;? This action
                            cannot be undone. This will permanently delete your project and remove your data from our
                            servers.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel disabled={fetcher.state === 'submitting' ? true : false}>
                            Cancel
                        </AlertDialogCancel>
                        <Button
                            variant="destructive"
                            onClick={handleProjectDeletion}
                            disabled={fetcher.state === 'submitting' ? true : false}
                        >
                            {fetcher.state === 'submitting' && <Loader2 className="animate-spin" />}
                            Delete
                        </Button>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    );
}
