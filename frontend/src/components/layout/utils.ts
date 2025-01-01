import type { LoaderFunction, LoaderFunctionArgs } from 'react-router-dom';

const loadLayout: LoaderFunction = (args: LoaderFunctionArgs) => {
    console.info('loadLayout ::', args);
    // TODO :: fetch the user info
    // TODO :: if we find no user loggedIn then redirect to login page via redirect command - https://reactrouter.com/6.28.1/fetch/redirect

    const data = {
        user: {
            name: 'dhyey',
            email: 'dhyey@example.com',
            iconUrl: '',
        },
    };
    return data;
};

export { loadLayout };
