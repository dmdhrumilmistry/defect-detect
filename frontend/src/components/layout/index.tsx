import { Outlet } from 'react-router-dom';

// shadcn/ui components
import { SidebarProvider, SidebarInset } from '@/components/ui/sidebar';

// internal components
import AppSidebar from '@/components/app-sidebar';
import Header from '@/components/header';

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
