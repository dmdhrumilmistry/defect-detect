import type { FetchArgs, Maybe, MutateArgs } from '@/types';

export default class RestService {
    static async fetch<T = unknown>(args: FetchArgs): Promise<Maybe<T>> {
        try {
            const response: Response = await fetch(args.url, {
                method: 'GET',
                headers: args.headers ?? {},
                signal: args.signal ?? null,
            });

            // reject 4xx/5xx errors
            if (!response.ok)
                throw new Error(`[ERR] GET action for ${args.url}; ${response.status}:${response.statusText}`);

            return (await response.json()) as T;
        } catch (err) {
            if (args.throwError) throw err;
            else {
                console.error(err);
                return null;
            }
        }
    }

    static async mutate<T = unknown>(args: MutateArgs): Promise<Maybe<T>> {
        try {
            const response: Response = await fetch(args.url, {
                method: args.method,
                body: args.body,
                headers: {
                    'Content-Type': 'application/json',
                    ...(args.headers ?? {}),
                },
                signal: args.signal ?? null,
            });

            // reject 4xx/5xx errors
            if (!response.ok)
                throw new Error(
                    `[ERR] ${args.method} action for ${args.url}; ${response.status}:${response.statusText}`
                );

            return (await response.json()) as T;
        } catch (err) {
            if (args.throwError) throw err;
            else {
                console.error(err);
                return null;
            }
        }
    }
}
