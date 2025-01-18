import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function sleep(sec: number) {
    return new Promise((resolve) => setTimeout(() => resolve('done'), sec * 1000));
}
