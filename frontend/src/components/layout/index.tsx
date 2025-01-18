import { Outlet } from 'react-router-dom';

import Header from '@/components/header';
import AppSidebar from '@/components/app-sidebar';
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';

export default function Layout() {
    console.info('[COMP] Layout');

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
