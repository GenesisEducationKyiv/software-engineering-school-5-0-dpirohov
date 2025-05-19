import * as React from 'react';
import {Outlet} from 'react-router';
import {DashboardLayout} from '@toolpad/core/DashboardLayout';
import {Box} from '@mui/material';

export default function Layout() {
    return (
        <DashboardLayout defaultSidebarCollapsed={true}>
            <Box sx={{padding: 5}}>
                <Outlet />
            </Box>
        </DashboardLayout>
    );
}
