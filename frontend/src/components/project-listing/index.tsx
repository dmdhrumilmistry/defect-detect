import { useLoaderData } from 'react-router-dom';

export default function ProjectListing() {
    const data = useLoaderData();
    console.info('[COMP] ProjectListing :: ', data);

    return (
        <>
            <h1>Project Listing</h1>
        </>
    );
}
