import {Dialog, DialogTitle, DialogContent, DialogActions, MenuItem, Button, CircularProgress, Box} from '@mui/material';
import {useForm, TextFieldElement, SelectElement} from 'react-hook-form-mui';
import {FC, useState} from 'react';
import {useNotifications} from '@toolpad/core';
import {SubscriptionService} from '../../api';

type CreateOrgDialogProps = {
    open: boolean;
    onClose: () => void;
};

type SubscribeForUpdatesForm = {
    city: string;
    email: string;
    frequency: string;
};

const SubscribeDialog: FC<CreateOrgDialogProps> = ({open, onClose}) => {
    const [pending, setPending] = useState(false);
    const notifications = useNotifications();
    const {
        reset,
        handleSubmit,
        control,
        formState: {errors, isSubmitting},
    } = useForm<SubscribeForUpdatesForm>({
        defaultValues: {
            city: '',
            email: '',
            frequency: 'daily',
        },
    });

    const onSubmit = async (data: SubscribeForUpdatesForm) => {
        try {
            setPending(true);
            await SubscriptionService.subscribe(data.email, data.city, data.frequency as 'daily' | 'hourly');
            notifications.show('Subscription submitted! Check your email for confirmation link!', {
                severity: 'success',
                autoHideDuration: 3000,
            });
        } catch (err: any) {
            console.error(err);
            notifications.show(err?.body?.message || 'Failed to subscribe', {severity: 'error', autoHideDuration: 3000});
        } finally {
            setPending(false);
            onClose();
        }
    };

    const frequencyData = [
        {label: 'Every hour', id: 'hourly'},
        {label: 'Every day', id: 'daily'},
    ];

    return (
        <Dialog open={open} onClose={onClose}>
            <DialogTitle>Subscribe for updates</DialogTitle>
            {pending ? (
                <DialogContent sx={{display: 'flex', justifyContent: 'center', p: 6}}>
                    <Box
                        sx={{
                            width: '100%',
                            height: '100%',
                            minWidth: 600,
                            minHeight: 250,
                            display: 'flex',
                            alignItems: 'center',
                            justifyContent: 'center',
                        }}
                    >
                        <CircularProgress />
                    </Box>
                </DialogContent>
            ) : (
                <form onSubmit={handleSubmit(onSubmit)} noValidate>
                    <DialogContent sx={{display: 'flex', flexDirection: 'column', gap: 2, minWidth: 600}}>
                        <TextFieldElement name="city" label={'City'} required fullWidth autoFocus control={control} />
                        <TextFieldElement name="email" label={'Email'} required fullWidth autoFocus control={control} />
                        <SelectElement
                            label={'Frequency'}
                            size="small"
                            control={control}
                            name="frequency"
                            options={frequencyData}
                            fullWidth
                            required={false}
                        />
                    </DialogContent>
                    <DialogActions>
                        <Button
                            onClick={() => {
                                onClose();
                                reset();
                            }}
                        >
                            Cancel
                        </Button>
                        <Button type="submit" variant="contained" disabled={false}>
                            Subscribe
                        </Button>
                    </DialogActions>
                </form>
            )}
        </Dialog>
    );
};

export default SubscribeDialog;
