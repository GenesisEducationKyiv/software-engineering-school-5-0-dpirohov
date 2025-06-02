import axios from 'axios';
import {QueryClient} from '@tanstack/react-query';

const protocol = window.location.protocol;
const URL = window.location.host;

axios.defaults.baseURL = `${protocol}//${URL}`;

export const queryClient = new QueryClient({
    defaultOptions: {
        queries: {
            refetchOnWindowFocus: false,
        },
    },
});
