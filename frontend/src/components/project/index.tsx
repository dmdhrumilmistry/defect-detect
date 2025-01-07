import type { ProjectLoader } from '@/types';
import { useLoaderData } from 'react-router-dom';

export default function Project() {
    const { project } = useLoaderData() as ProjectLoader;
    console.info('[COMP] Project :: ', project);

    if (!project) return <div>Project failed to fetch!</div>;
    return (
        <>
            <h1>{project.title}</h1>
        </>
    );
}
