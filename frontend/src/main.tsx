import * as React from 'react';
import * as ReactDOM from 'react-dom/client';
import {createBrowserRouter, RouterProvider} from 'react-router';
import App from './App';
import Layout from './layouts/dashboard';
import DashboardPage from './pages';
import ConfirmPage from './pages/ConfirmationPage/ConfirmationPage';
import NotFound from './pages/NotFound/NotFound';
import UnsubscribePage from './pages/UnsubscribePage/UnsubscribePage';

const router = createBrowserRouter([
    {
        Component: App,
        children: [
            {
                path: '/',
                Component: Layout,
                children: [
                    {
                        path: '',
                        Component: DashboardPage,
                    },
                ],
            },
            {
                path: '/confirm/:token',
                element: <ConfirmPage />,
            },
            {
                path: '/unsubscribe/:token',
                element: <UnsubscribePage />,
            },
            {
                path: '*',
                Component: NotFound,
            },
        ],
    },
]);

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <RouterProvider router={router} />
    </React.StrictMode>,
);
