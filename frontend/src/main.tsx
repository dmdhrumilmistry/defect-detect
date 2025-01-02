import './global.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

// UI Components
import Layout from '@/components/layout';
import LayoutErrorBoundary from '@/components/layout/error-boundary';
import ProjectListing from '@/components/project-listing';
import Project from '@/components/project';
import ProjectErrorBoundary from '@/components/project/error-boundary';

// data loaders
import { loadLayout } from '@/components/layout/utils';
import { loadProjects } from '@/components/project-listing/utils';
import { loadProject } from '@/components/project/utils';

const router = createBrowserRouter([
    {
        id: 'layout',
        path: '/',
        loader: loadLayout,
        element: <Layout />,
        hasErrorBoundary: true,
        errorElement: <LayoutErrorBoundary />,
        children: [
            {
                id: 'dashboard',
                index: true,
                element: <div>Dashboard</div>,
            },
            {
                id: 'projectListing',
                path: '/projects',
                loader: loadProjects,
                element: <ProjectListing />,
            },
            {
                id: 'project',
                path: '/projects/:projectId',
                loader: loadProject,
                element: <Project />,
                hasErrorBoundary: true,
                errorElement: <ProjectErrorBoundary />,
            },
        ],
    },
]);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} fallbackElement={<div>Parent loader...</div>} />
    </StrictMode>
);

/**
 * Note ::
 * On page load createBrowserRouter will initiate all matching route loaders when it mounts.
 * During this time fallbackElement will be rendered (if provided) to indicate app is working in background.
 */
