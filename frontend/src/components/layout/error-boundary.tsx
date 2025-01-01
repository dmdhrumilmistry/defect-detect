import { useRouteError } from 'react-router-dom';

export default function LayoutErrorBoundary() {
    const error = useRouteError();
    console.error('LayoutErrorBoundary :: ', error);

    // TODO :: handle this UI
    return (
        <>
            <div>Layout Error Boundary</div>
        </>
    );
}
