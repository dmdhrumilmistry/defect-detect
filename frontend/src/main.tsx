import './styles/global.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';

// UI Components
import Layout from '@/components/layout';
import LayoutSkeleton from './components/layout/skeleton';
import LayoutErrorBoundary from '@/components/layout/error-boundary';
import Projects from '@/components/projects';
import Project from '@/components/project';
import ProjectErrorBoundary from '@/components/project/error-boundary';

// data loaders
import { layoutDataLoader, layoutHandle } from '@/components/layout/utils';
import { projectsDataLoader, projectsAction, projectsHandle } from '@/components/projects/utils';
import { projectDataLoader, projectHandle, projectAction } from '@/components/project/utils';

const router = createBrowserRouter([
    {
        id: 'layoutRoute',
        element: <Layout />,
        errorElement: <LayoutErrorBoundary />,
        loader: layoutDataLoader,
        handle: layoutHandle,
        children: [
            {
                id: 'dashboardRoute',
                index: true,
                element: <div>Dashboard</div>,
                // no breadcrumb as its the index handler
            },
            {
                id: 'projectsRoute',
                path: '/projects',
                handle: projectsHandle,
                children: [
                    {
                        id: 'allProjectsRoute',
                        index: true,
                        element: <Projects />,
                        loader: projectsDataLoader,
                        action: projectsAction,
                        // no breadcrumb as its the index handler
                    },
                    {
                        id: 'projectRoute',
                        path: '/projects/:projectId',
                        element: <Project />,
                        errorElement: <ProjectErrorBoundary />,
                        loader: projectDataLoader,
                        action: projectAction,
                        handle: projectHandle,
                    },
                ],
            },
        ],
    },
]);

createRoot(document.getElementById('root')!).render(
    <StrictMode>
        <RouterProvider router={router} fallbackElement={<LayoutSkeleton />} />
    </StrictMode>
);

/**
 * Note ::
 * On page load createBrowserRouter will initiate all matching route loaders when it mounts.
 * During this time fallbackElement will be rendered (if provided) to indicate app is working in background.
 */
