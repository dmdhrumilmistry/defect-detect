import type { TUser } from '@/types';
import { Link, Outlet, useLoaderData, useNavigation } from 'react-router-dom';

export default function Layout() {
    const { user } = useLoaderData() as { user: TUser };
    const navigation = useNavigation();
    console.info('[COMP] Layout :: ', user, navigation);

    return (
        <>
            {/* TODO :: Handle the surrounding layout UI */}
            <h1>This is layout</h1>
            <code>UserName: {user.username}</code>
            <br />
            <Link to="/projects">Navigate &gt; Projects</Link>
            <br />
            <Link to="/projects/30">Navigate &gt; Test Project</Link>
            <Outlet />
        </>
    );
}
