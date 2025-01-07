import type { LayoutLoader } from '@/types';
import { Outlet, useLoaderData, useNavigation } from 'react-router-dom';

import Header from '@/components/header';
import AppSidebar from '@/components/app-sidebar';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';

export default function Layout() {
    const { user } = useLoaderData() as LayoutLoader;
    const navigation = useNavigation();
    console.info('[COMP] Layout :: ', user, navigation);

    return (
        <SidebarProvider>
            <AppSidebar />
            <SidebarInset>
                <Header />
                <Outlet />
            </SidebarInset>
        </SidebarProvider>
    );
}
