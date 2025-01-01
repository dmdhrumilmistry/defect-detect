import { Outlet, useLoaderData, useNavigation } from 'react-router-dom';

export default function Layout() {
    const data = useLoaderData();
    const navigation = useNavigation();
    console.info('Layout :: ', data, navigation);

    return (
        <>
            {/* TODO :: Handle the surrounding layout UI */}
            <h1>This is layout</h1>
            <Outlet />
        </>
    );
}
