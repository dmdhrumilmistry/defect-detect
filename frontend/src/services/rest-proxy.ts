import type { FetchProxyArgs, Maybe } from '../types';

import RestService from './rest';

export default class RestServiceProxy {
    static map = new Map();

    static async fetch<T = unknown>(args: FetchProxyArgs): Promise<Maybe<T>> {
        const { cacheKey, ...fetchArgs } = args;

        // cache lookup
        if (this.map.has(cacheKey)) return this.map.get(cacheKey) as T;

        // cache miss
        const data = await RestService.fetch<T>(fetchArgs);
        if (!data) return data;

        // cache save
        this.map.set(cacheKey, data);
        return data;
    }

    static invalidateCache(cacheKey: string): boolean {
        return this.map.delete(cacheKey);
    }
}
