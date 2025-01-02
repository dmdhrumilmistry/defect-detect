import type { Maybe } from '@/types';

export default class CookieService {
    static get(name: string): Maybe<string> {
        const match = new RegExp('(^|;)\\s*' + encodeURIComponent(name) + '=([^;]*)').exec(document.cookie);
        return match ? decodeURIComponent(match[2]) : null;
    }

    static set(name: string, value: string, maxAgeInSec?: number): void {
        const defaultCookieSettings = 'secure=true; samesite=Strict; path=/';
        let cookieString = `${encodeURIComponent(name)}=${encodeURIComponent(value)}; ${defaultCookieSettings}`;
        if (maxAgeInSec) cookieString += `; max-age=${maxAgeInSec}`;
        document.cookie = cookieString;
    }

    static delete(name: string): void {
        this.set(name, '', -1);
    }
}