import CloudIcon from '@mui/icons-material/Cloud';
import {Outlet} from 'react-router';
import {ReactRouterAppProvider} from '@toolpad/core/react-router';
import type {Navigation} from '@toolpad/core/AppProvider';
import {NotificationsProvider} from '@toolpad/core';
import {QueryClientProvider} from '@tanstack/react-query';
import {queryClient} from './api/context/apiClient';

const NAVIGATION: Navigation = [
    {
        title: 'Weather',
        icon: <CloudIcon />,
    },

];

const BRANDING = {
    title: 'WeatherApi',
};

export default function App() {
    return (
        <NotificationsProvider>
            <QueryClientProvider client={queryClient}>
                <ReactRouterAppProvider navigation={NAVIGATION} branding={BRANDING}>
                    <Outlet />
                </ReactRouterAppProvider>
            </QueryClientProvider>
        </NotificationsProvider>
    );
}
