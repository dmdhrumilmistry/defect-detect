import type { LoaderFunction, LoaderFunctionArgs, ActionFunction, ActionFunctionArgs } from 'react-router-dom';
import type { LoaderReturnValue, Maybe, ProjectLoader, RouteHandle, StringKVs, TProject } from '@/types';
import { API_BASE_URL, CACHE_KEYS } from '@/services/const';
import RestServiceProxy from '@/services/rest-proxy';
import { sleep } from '@/lib/utils';

const projectLoader: LoaderFunction = async (args: LoaderFunctionArgs): Promise<ProjectLoader> => {
    console.info('[LOADER] Project :: ', args);
    const projectId = args.params.projectId;
    if (!projectId) throw new Error('Project Id not found!');

    const project = await RestServiceProxy.fetch<TProject>({
        url: `${API_BASE_URL}/products/${projectId}`,
        cacheKey: `${CACHE_KEYS.projectPrefix}${projectId}`,
        signal: args.request.signal,
        throwError: true,
    });
    if (!project) throw new Error('Failed to fetch the project!');
    return { project };
};

const projectAction: ActionFunction = async (args: ActionFunctionArgs): Promise<StringKVs> => {
    console.log('[ACTION] Project :: ', args);
    const projectId = args.params.projectId;
    if (!projectId) throw new Error('Project Id not found!');

    if (args.request.method === 'DELETE') {
        // TODO :: handle project deletion
    } else if (args.request.method === 'PUT') {
        // TODO :: handle project updates
    } else {
        throw new Error('Unsupported action method!');
    }

    await sleep(5);
    // invalidate cache
    RestServiceProxy.invalidateCache(CACHE_KEYS.projects);
    RestServiceProxy.invalidateCache(`${CACHE_KEYS.projectPrefix}${projectId}`);
    return { success: 'true' };
};

const projectHandle: RouteHandle = {
    breadcrumb: (data: Maybe<LoaderReturnValue>) => {
        if (data && 'project' in data) return { href: `/projects/${data.project.id}`, label: data.project.title };
        else return { href: '/projects/null', label: 'Project' };
    },
};

export { projectLoader, projectAction, projectHandle };
