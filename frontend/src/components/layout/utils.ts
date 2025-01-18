import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';
import type { LayoutDataLoader, RouteHandle, TUser } from '@/types';
import { COOKIE_KEYS, API_BASE_URL, CACHE_KEYS } from '@/services/const';
import RestServiceProxy from '@/services/rest-proxy';
import CookieService from '@/services/cookie';

const layoutDataLoader: LoaderFunction = async (args: LoaderFunctionArgs): Promise<LayoutDataLoader> => {
    console.info('[LOADER] Layout ::', args);
    const userId = CookieService.get(COOKIE_KEYS.userId);
    if (!userId) throw new Error('User Id not found!');

    const user = await RestServiceProxy.fetch<TUser>({
        url: `${API_BASE_URL}/users/${userId}`,
        cacheKey: CACHE_KEYS.user,
        signal: args.request.signal,
        throwError: true,
    });
    if (!user) throw new Error('User undefined');
    return { user };
};

const layoutHandle: RouteHandle = {
    breadcrumb: () => ({ href: '/', label: 'Dashboard' }),
};

export { layoutDataLoader, layoutHandle };

// TODO :: if we find no loggedIn user then redirect to login page via redirect command - https://reactrouter.com/6.28.1/fetch/redirect
