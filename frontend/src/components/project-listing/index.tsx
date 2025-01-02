import type { TProject } from '@/types';
import { useLoaderData, useNavigation } from 'react-router-dom';

export default function ProjectListing() {
    const { projects } = useLoaderData() as { projects: TProject[] };
    const navigation = useNavigation();
    console.info('[COMP] ProjectListing :: ', projects, navigation);

    return (
        <>
            <h1>Project Listing</h1>
            <code>Total projects: {projects.length}</code>
        </>
    );
}
