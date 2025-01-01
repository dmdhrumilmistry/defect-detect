import './global.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

// UI Components
import Layout from './components/layout';
import LayoutErrorBoundary from './components/layout/error-boundary';
import ProjectListing from './components/project-listing';
import Project from './components/project';

// data loaders
import { loadLayout } from './components/layout/utils';
import { loadProjects } from './components/project-listing/utils';
import { loadProject } from './components/project/utils';

const router = createBrowserRouter([
    {
        id: 'layout',
        path: '/',
        loader: loadLayout,
        hasErrorBoundary: true,
        errorElement: <LayoutErrorBoundary />,
        element: <Layout />,
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
            },
        ],
    },
]);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} fallbackElement={<div>Parent loader...</div>} />
    </StrictMode>
);
