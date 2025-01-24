// types
import type { ProjectsActionError } from '@/types';
import type { ChangeEvent } from 'react';

// libs
import { useEffect } from 'react';
import { z } from 'zod';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { useFetcher } from 'react-router-dom';
import { Plus, Loader2, AlertCircle } from 'lucide-react';

// shadcn/ui components
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group';
import { Alert, AlertTitle, AlertDescription } from '@/components/ui/alert';
import { Dialog, DialogTrigger, DialogContent, DialogHeader, DialogFooter, DialogTitle } from '@/components/ui/dialog';
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';

// utilities
import { convertFileToString } from '@/lib/utils';
import { subtleText } from '@/styles/standard-classes';

const modeEnum = z.enum(['url', 'file']);
const createProjectFormSchema = z.discriminatedUnion('mode', [
    z.object({
        mode: modeEnum.extract(['url']),
        projectName: z.string().min(3, 'Project name must be at least 3 characters.'),
        repoUrl: z.string().url('Provide valid repository url.'),
    }),
    z.object({
        mode: modeEnum.extract(['file']),
        projectName: z.string().min(3, 'Project name must be at least 3 characters.'),
        sbomJsonFile: z.string({
            required_error: 'Select valid JSON files.',
        }),
    }),
]);
export type CreateProjectFormSchema = z.infer<typeof createProjectFormSchema>;

export default function CreateProject() {
    const form = useForm<CreateProjectFormSchema>({
        resolver: zodResolver(createProjectFormSchema),
        defaultValues: {
            mode: 'url',
            projectName: '',
            repoUrl: '',
        },
    });
    const { unregister } = form;
    const mode = form.getValues('mode');

    const fetcher = useFetcher<ProjectsActionError>();
    const isSubmitting = fetcher.state === 'submitting';
    console.log('[COMP] CreateProject :: ', form, fetcher);

    useEffect(() => {
        const unregisterField = mode === 'url' ? 'sbomJsonFile' : 'repoUrl';
        unregister(unregisterField, {
            keepValue: false,
            keepError: false,
            keepDefaultValue: true,
        });
    }, [mode, unregister]);

    const handleProjectCreation = (data: CreateProjectFormSchema) => {
        fetcher.submit(data, {
            method: 'post',
            action: '/projects?index',
            encType: 'application/json',
        });
    };
    const handleFileInputOnChange = (event: ChangeEvent<HTMLInputElement>) => {
        const files = event.target?.files;
        console.log(files);
        if (!files?.length) return;

        convertFileToString(files[0])
            .then((sbomJsonFile) => {
                form.setValue('sbomJsonFile', sbomJsonFile);
                form.clearErrors('sbomJsonFile');
            })
            .catch((err) => {
                form.resetField('sbomJsonFile');
                console.error(err);
            });
    };

    return (
        <Dialog>
            <DialogTrigger asChild>
                <Button>
                    <span className="hidden sm:block">New Project</span>
                    <Plus className="block sm:hidden !size-5" />
                </Button>
            </DialogTrigger>
            <DialogContent className="max-w-[80%] sm:max-w-lg gap-0" disableClose={isSubmitting}>
                <DialogHeader className="mb-6">
                    <DialogTitle>Create a new project</DialogTitle>
                </DialogHeader>

                {fetcher.data?.error && (
                    <Alert variant="destructive" className="mb-4">
                        <AlertCircle className="h-4 w-4" />
                        <AlertTitle>{fetcher.data?.error?.status}</AlertTitle>
                        <AlertDescription>{fetcher.data?.error?.message}</AlertDescription>
                    </Alert>
                )}

                <Form {...form}>
                    {/* eslint-disable-next-line @typescript-eslint/no-misused-promises */}
                    <form onSubmit={form.handleSubmit(handleProjectCreation)} className="flex flex-col gap-4">
                        <FormField
                            control={form.control}
                            name="projectName"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel className={`block ${subtleText}`}>Project Name</FormLabel>
                                    <FormControl>
                                        <Input
                                            {...field}
                                            autoComplete="off"
                                            placeholder="Proj..."
                                            disabled={isSubmitting}
                                        />
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="mode"
                            render={({ field }) => (
                                <FormItem>
                                    <FormLabel className={`block ${subtleText}`}>Mode of SBOM</FormLabel>
                                    <FormControl>
                                        <RadioGroup
                                            className="flex flex-row gap-6"
                                            orientation="horizontal"
                                            onValueChange={field.onChange}
                                            defaultValue={field.value}
                                            disabled={isSubmitting}
                                        >
                                            <FormItem className="flex items-center gap-2 space-y-0">
                                                <FormControl>
                                                    <RadioGroupItem value="url" />
                                                </FormControl>
                                                <FormLabel className="font-normal cursor-pointer">Repo URL</FormLabel>
                                            </FormItem>
                                            <FormItem className="flex items-center gap-2 space-y-0">
                                                <FormControl>
                                                    <RadioGroupItem value="file" />
                                                </FormControl>
                                                <FormLabel className="font-normal cursor-pointer">JSON File</FormLabel>
                                            </FormItem>
                                        </RadioGroup>
                                    </FormControl>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        {mode === 'url' && (
                            <FormField
                                control={form.control}
                                name="repoUrl"
                                render={({ field }) => (
                                    <FormItem>
                                        <FormLabel className={`block ${subtleText}`}>Repository URL</FormLabel>
                                        <FormControl>
                                            <Input
                                                {...field}
                                                autoComplete="off"
                                                placeholder="https://..."
                                                disabled={isSubmitting}
                                            />
                                        </FormControl>
                                        <FormDescription>Extract SBOM content using repo.</FormDescription>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        )}
                        {mode === 'file' && (
                            <FormField
                                control={form.control}
                                name="sbomJsonFile"
                                render={() => (
                                    <FormItem>
                                        <FormLabel className={`block ${subtleText}`}>SBOM JSON file</FormLabel>
                                        <FormControl>
                                            <Input
                                                type="file"
                                                accept="application/json"
                                                onChange={(event) => handleFileInputOnChange(event)}
                                                disabled={isSubmitting}
                                            />
                                        </FormControl>
                                        <FormMessage />
                                    </FormItem>
                                )}
                            />
                        )}
                        <DialogFooter>
                            <Button type="submit" disabled={isSubmitting}>
                                {isSubmitting && <Loader2 className="animate-spin" />}Create
                            </Button>
                        </DialogFooter>
                    </form>
                </Form>
            </DialogContent>
        </Dialog>
    );
}
