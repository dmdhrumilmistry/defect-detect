import { useRouteError } from 'react-router-dom';

export default function ProjectErrorBoundary() {
    const error = useRouteError();
    console.error('[COMP] ProjectErrorBoundary :: ', error);

    // TODO :: handle this UI
    return (
        <>
            <div>Project Error Boundary</div>
        </>
    );
}
