import {useEffect, useState} from 'react';
import {Box, CircularProgress, Typography, Link as MuiLink} from '@mui/material';
import {useNotifications} from '@toolpad/core';
import {SubscriptionService} from '../../api';
import {useParams, Link} from 'react-router';
import {CONFIRM_PAGE_IDS} from '../../constants/test_ids';

export default function ConfirmPage() {
    const {token} = useParams<{token: string}>();
    const notifications = useNotifications();
    const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');

    useEffect(() => {
        if (!token) return;

        SubscriptionService.confirmSubscription(token)
            .then(() => {
                setStatus('success');
                notifications.show('Subscription confirmed!', {severity: 'success', autoHideDuration: 3000});
            })
            .catch(err => {
                setStatus('error');
                notifications.show(`Confirmation failed: ${err.message}`, {severity: 'error', autoHideDuration: 3000});
            });
    }, [token]);

    return (
        <Box display="flex" justifyContent="center" alignItems="center" minHeight="80vh" flexDirection="column">
            {status === 'loading' && <CircularProgress />}
            {status === 'success' && (
                <Typography variant="h5" data-testid={CONFIRM_PAGE_IDS.confirmation}>
                    Subscription confirmed ✅
                </Typography>
            )}
            {status === 'error' && (
                <Typography variant="h5" color="error" data-testid={CONFIRM_PAGE_IDS.confirmation}>
                    Confirmation failed ❌
                </Typography>
            )}
            <MuiLink component={Link} to="/" underline="hover" data-testid={CONFIRM_PAGE_IDS.linkToMainPage}>
                Back to main page
            </MuiLink>
        </Box>
    );
}
