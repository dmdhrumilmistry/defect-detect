export type Maybe<T> = T | null | undefined;

export type StringKVs = Record<string, string>;

export type HttpMethod = 'PUT' | 'POST';

export type FetchArgs = Readonly<{
    url: string;
    headers?: HeadersInit;
    signal?: AbortSignal;
    throwError?: boolean;
}>;

export type FetchProxyArgs = FetchArgs &
    Readonly<{
        cacheKey: string;
    }>;

export type MutateArgs = Readonly<{
    method: HttpMethod;
    url: string;
    body: string;
    headers?: HeadersInit;
    signal?: AbortSignal;
    throwError?: boolean;
}>;

export type TUser = Readonly<{
    id: string;
    firstName: string;
    lastName: string;
    email: string;
    username: string;
    image: string;
}>;

export type TProject = Readonly<{
    id: string;
    title: string;
    description: string;
}>;
