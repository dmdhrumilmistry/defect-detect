import { useLoaderData } from 'react-router-dom';
import type { StringKVs } from '../../types';

export default function Project() {
    const data = useLoaderData() as StringKVs;
    console.info('[COMP] Project :: ', data);

    return (
        <>
            <h1>{data.name}</h1>
        </>
    );
}
