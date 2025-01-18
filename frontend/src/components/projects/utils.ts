import type { LoaderFunction, LoaderFunctionArgs, ActionFunction, ActionFunctionArgs } from 'react-router-dom';
import type { ProjectsDataLoader, RouteHandle, StringKVs, TProject } from '@/types';
import { API_BASE_URL, CACHE_KEYS } from '@/services/const';
import RestServiceProxy from '@/services/rest-proxy';
import { sleep } from '@/lib/utils';

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

const projectsAction: ActionFunction = async (args: ActionFunctionArgs): Promise<StringKVs> => {
    console.log('[ACTION] Project(s) :: ', args);

    if (args.request.method === 'POST') {
        // TODO :: handle project creation
    } else {
        throw new Error('Unsupported action method!');
    }

    await sleep(5);
    // invalidate cache
    RestServiceProxy.invalidateCache(CACHE_KEYS.projects);
    return { success: 'true' };
};

const projectsHandle: RouteHandle = {
    breadcrumb: () => ({ href: '/projects', label: 'Projects' }),
};

export { projectsDataLoader, projectsAction, projectsHandle };
