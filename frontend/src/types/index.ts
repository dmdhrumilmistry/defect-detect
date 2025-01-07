export type Maybe<T> = T | null | undefined;

export type StringKVs = Record<string, string>;

export type HttpMethod = 'PUT' | 'POST';

export type FetchArgs = {
    readonly url: string;
    readonly headers?: HeadersInit;
    readonly signal?: AbortSignal;
    readonly throwError?: boolean;
};

export type FetchProxyArgs = FetchArgs & {
    readonly cacheKey: string;
};

export type MutateArgs = {
    readonly method: HttpMethod;
    readonly url: string;
    readonly body: string;
    readonly headers?: HeadersInit;
    readonly signal?: AbortSignal;
    readonly throwError?: boolean;
};

export type TUser = {
    readonly id: string;
    readonly firstName: string;
    readonly lastName: string;
    readonly email: string;
    readonly username: string;
    readonly image: string;
};

export type TProject = {
    readonly id: string;
    readonly title: string;
    readonly description: string;
};

// Data loader function return type
export type LayoutLoader = {
    readonly user: TUser;
};

export type ProjectsLoader = {
    readonly projects: TProject[];
};

export type ProjectLoader = {
    readonly project: TProject;
};

export type LoaderReturnValue = LayoutLoader | ProjectsLoader | ProjectLoader;
// -----

export type RouteHandle = {
    readonly breadcrumb?: (data?: Maybe<LoaderReturnValue>) => { href: string; label: string };
};
