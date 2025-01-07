import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';
import type { LoaderReturnValue, Maybe, ProjectLoader, RouteHandle, TProject } from '@/types';
import { API_BASE_URL, CACHE_KEYS } from '@/services/const';
import RestServiceProxy from '@/services/rest-proxy';

const projectLoader: LoaderFunction = async (args: LoaderFunctionArgs): Promise<ProjectLoader> => {
    console.info('[LOADER] Project :: ', args);
    if (!args.params.projectId) throw new Error('Project Id not found!');

    const project = await RestServiceProxy.fetch<TProject>({
        url: `${API_BASE_URL}/products/${args.params.projectId}`,
        cacheKey: `${CACHE_KEYS.projectPrefix}${args.params.projectId}`,
        signal: args.request.signal,
        throwError: true,
    });
    if (!project) throw new Error('Failed to fetch the project!');
    return { project };
};

const projectHandle: RouteHandle = {
    breadcrumb: (data: Maybe<LoaderReturnValue>) => {
        if (data && 'project' in data) return { href: `/projects/${data.project.id}`, label: data.project.title };
        else return { href: '/projects/null', label: 'Project' };
    },
};

export { projectLoader, projectHandle };
