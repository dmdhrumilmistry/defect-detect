import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs));
}

export function sleep(sec: number) {
    return new Promise((resolve) => setTimeout(() => resolve('done'), sec * 1000));
}

export async function convertFileToString(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
        try {
            const reader = new FileReader();
            reader.onload = (event) => {
                resolve(event.target?.result as string);
            };
            reader.readAsText(file);
        } catch (err) {
            console.error(err);
            reject(new Error('Error while parsing File'));
        }
    });
}

export function convertStringToFile(content: string, fileName: string): File {
    return new File([content], fileName, {
        type: 'application/json',
    });
}
