import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';

const loadProject: LoaderFunction = (args: LoaderFunctionArgs) => {
    // TODO :: get the project ID
    console.info('[LOADER] Project :: ', args.params.projectId);

    return {
        name: 'Dummy Project 1',
        id: 'dummy1',
    };
};

export { loadProject };
