import {FC, memo} from 'react';
import {Box, Card, CardContent, Skeleton, Stack, Typography} from '@mui/material';

const LoadingSkeleton: FC = () => {
    return (
        <Stack direction="column" spacing={3}>
            <Skeleton variant="rounded" height={40} />
            <Stack direction="row" spacing={2}>
                <Card variant="outlined" sx={{flex: 1}}>
                    <CardContent>
                        <Typography variant="caption">
                            <Skeleton width={80} />
                        </Typography>
                        <Typography variant="h4">
                            <Skeleton width={100} />
                        </Typography>
                    </CardContent>
                </Card>
                <Card variant="outlined" sx={{flex: 1}}>
                    <CardContent>
                        <Typography variant="caption">
                            <Skeleton width={80} />
                        </Typography>
                        <Typography variant="h4">
                            <Skeleton width={100} />
                        </Typography>
                    </CardContent>
                </Card>
            </Stack>
            <Card variant="outlined">
                <CardContent>
                    <Typography variant="caption">
                        <Skeleton width={120} />
                    </Typography>
                    <Typography variant="h4">
                        <Skeleton width="80%" />
                    </Typography>
                </CardContent>
            </Card>
            <Box display="flex" justifyContent="center" mt={2}>
                <Skeleton variant="rectangular" width={200} height={36} />
            </Box>
        </Stack>
    );
};

export default memo(LoadingSkeleton);
