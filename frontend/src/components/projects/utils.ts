import {
    type LoaderFunction,
    type LoaderFunctionArgs,
    type ActionFunction,
    type ActionFunctionArgs,
    redirect,
} from 'react-router-dom';
import type { ProjectsDataLoader, RouteHandle, TProject } from '@/types';
import type { CreateProjectFormSchema } from './create-project';

import RestServiceProxy from '@/services/rest-proxy';
import { API_BASE_URL, CACHE_KEYS } from '@/services/const';
import { convertStringToFile, sleep } from '@/lib/utils';

const projectsDataLoader: LoaderFunction = async (args: LoaderFunctionArgs): Promise<ProjectsDataLoader> => {
    console.info('[LOADER] Project(s) ::', args);

    const projects = await RestServiceProxy.fetch<unknown>({
        url: `${API_BASE_URL}/products`,
        cacheKey: CACHE_KEYS.projects,
        signal: args.request.signal,
        throwError: true,
    });

    // @ts-expect-error ignore the below error for now
    return { projects: (projects?.products ?? []) as TProject[] };
};

const projectsAction: ActionFunction = async (args: ActionFunctionArgs) => {
    const data = (await args.request.json()) as CreateProjectFormSchema;
    console.log('[ACTION] Project(s) :: ', args, data);

    if (args.request.method === 'POST') {
        // TODO :: handle project creation
        if (data.mode === 'file') console.log(convertStringToFile(data.sbomJsonFile, 'file.json'));
        await sleep(2);
    } else {
        throw new Error('Unsupported action method!');
    }

    // invalidate cache
    RestServiceProxy.invalidateCache(CACHE_KEYS.projects);
    return redirect('/projects/20');
    // simulate the endpoint error
    // return {
    //     error: {
    //         status: 'Failed',
    //         message: 'Oops something went wrong while create project.',
    //         code: 500,
    //     },
    // };
};

const projectsHandle: RouteHandle = {
    breadcrumb: () => ({ href: '/projects', label: 'Projects' }),
};

export { projectsDataLoader, projectsAction, projectsHandle };
