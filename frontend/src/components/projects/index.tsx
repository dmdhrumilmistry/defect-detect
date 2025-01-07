import type { ProjectsLoader } from '@/types';
import { useLoaderData, useNavigation } from 'react-router-dom';

export default function Projects() {
    const { projects } = useLoaderData() as ProjectsLoader;
    const navigation = useNavigation();
    console.info('[COMP] Projects :: ', projects, navigation);

    return (
        <>
            <h1>Project Listing</h1>
            <code>Total projects: {projects.length}</code>
        </>
    );
}
