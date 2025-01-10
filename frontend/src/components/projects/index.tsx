import type { ProjectsLoader, TProject } from '@/types';
import { useState } from 'react';
import { useLoaderData, useNavigation } from 'react-router-dom';
import { Plus } from 'lucide-react';

import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogFooter,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogCancel,
    AlertDialogAction,
} from '@/components/ui/alert-dialog';
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import SearchBar from '@/components/shared/search-bar';
import ProjectCard from '@/components/shared/project-card';
import { subtleText } from '@/styles/standard-classes';

export default function Projects() {
    const { projects } = useLoaderData() as ProjectsLoader;
    const navigation = useNavigation();
    console.info('[COMP] Projects :: ', projects, navigation);

    const [searchQuery, setSearchQuery] = useState('');
    const [projectMarkedForDeletion, setProjectMarkedForDeletion] = useState<TProject | null>(null);
    const filteredProjects = projects.filter((project) =>
        project.title.toLowerCase().includes(searchQuery.toLowerCase())
    );

    const clearProjectMarkedForDeletion = () => setProjectMarkedForDeletion(null);
    const handleProjectDeletion = () => {
        console.log('handleProjectDeletion :: ', projectMarkedForDeletion);
        clearProjectMarkedForDeletion();
        // TODO :: handle the project deletion action
    };

    return (
        <>
            <div className="min-h-screen mx-auto max-w-[1392px] w-full p-4">
                <div className="flex items-center justify-end gap-4">
                    <SearchBar searchQuery={searchQuery} setSearchQuery={setSearchQuery} />
                    <Dialog>
                        <DialogTrigger asChild>
                            <Button>
                                <span className="hidden sm:block">New Project</span>
                                <Plus className="block sm:hidden !size-5" />
                            </Button>
                        </DialogTrigger>
                        <DialogContent className="max-w-[80%] sm:max-w-lg gap-6">
                            <DialogHeader>
                                <DialogTitle>Create a new project</DialogTitle>
                            </DialogHeader>
                            <div className="flex gap-2 flex-col">
                                <Label htmlFor="name" className={subtleText}>
                                    Name<span className="text-red-500">*</span>
                                </Label>
                                <Input
                                    id="name"
                                    autoComplete="off"
                                    placeholder="Project name..."
                                    className="col-span-3"
                                />
                            </div>
                            <DialogFooter>
                                {/* TODO :: handle submission */}
                                <Button type="submit">Create</Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
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
            <AlertDialog open={Boolean(projectMarkedForDeletion)}>
                <AlertDialogContent className="max-w-[80%] sm:max-w-lg">
                    <AlertDialogHeader>
                        <AlertDialogTitle>{projectMarkedForDeletion?.title}</AlertDialogTitle>
                        <AlertDialogDescription>
                            Are you absolutely sure? This action cannot be undone. This will permanently delete your
                            project and remove your data from our servers.
                        </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                        <AlertDialogCancel onClick={clearProjectMarkedForDeletion}>Cancel</AlertDialogCancel>
                        <AlertDialogAction onClick={handleProjectDeletion}>Delete Project</AlertDialogAction>
                    </AlertDialogFooter>
                </AlertDialogContent>
            </AlertDialog>
        </>
    );
}
