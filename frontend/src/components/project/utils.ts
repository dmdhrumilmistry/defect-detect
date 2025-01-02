import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';
import type { TProject } from '../../types';
import { API_BASE_URL, CACHE_KEYS } from '../../services/const';
import RestServiceProxy from '../../services/rest-proxy';

const loadProject: LoaderFunction = async (args: LoaderFunctionArgs) => {
    console.info('[LOADER] Project :: ', args);
    if (!args.params.projectId) throw new Error('Project Id not found!');

    const project = await RestServiceProxy.fetch<TProject>({
        url: `${API_BASE_URL}/products/${args.params.projectId}`,
        cacheKey: `${CACHE_KEYS.projectPrefix}${args.params.projectId}`,
        signal: args.request.signal,
        throwError: true,
    });
    return { project };
};

export { loadProject };
