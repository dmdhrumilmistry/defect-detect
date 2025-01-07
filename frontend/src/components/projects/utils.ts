import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';
import type { ProjectsLoader, RouteHandle, TProject } from '@/types';
import { API_BASE_URL, CACHE_KEYS } from '@/services/const';
import RestServiceProxy from '@/services/rest-proxy';

const projectsLoader: LoaderFunction = async (args: LoaderFunctionArgs): Promise<ProjectsLoader> => {
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

const projectsHandle: RouteHandle = {
    breadcrumb: () => ({ href: '/projects', label: 'Projects' }),
};

export { projectsLoader, projectsHandle };
