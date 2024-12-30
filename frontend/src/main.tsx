import './global.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

import Layout from './components/layout';
import LayoutErrorBoundary from './components/layout/error-boundary';
import ProjectListing from './components/project-listing/page';

const router = createBrowserRouter([
    {
        id: 'layout',
        // loader: layoutLoader,
        hasErrorBoundary: true,
        errorElement: <LayoutErrorBoundary />,
        element: <Layout />,
        children: [
            {
                id: 'home',
                index: true,
                element: <div>Home</div>,
            },
            {
                id: 'projectListing',
                path: '/projects',
                // loader: projectLoader,
                element: <ProjectListing />,
            },
        ],
    },
]);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} fallbackElement={<div>Parent loader...</div>} />
    </StrictMode>
);
