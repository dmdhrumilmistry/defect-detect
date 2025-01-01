import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';

const loadProjects: LoaderFunction = (args: LoaderFunctionArgs) => {
    // TODO :: get the user info
    console.info('[LOADER] Project(s) ::', args);

    const data = [
        {
            name: 'Dummy Project 1',
            id: 'dummy1',
        },
        {
            name: 'Dummy Project 2',
            id: 'dummy2',
        },
    ];
    return data;
};

export { loadProjects };
