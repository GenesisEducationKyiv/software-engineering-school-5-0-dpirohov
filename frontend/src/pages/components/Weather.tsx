import {FC, memo, useState} from 'react';
import {Weather} from '../../api';
import {TextFieldElement, useForm} from 'react-hook-form-mui';
import {Box, Button, Card, CardContent, IconButton, InputAdornment, Stack, Typography} from '@mui/material';
import SearchIcon from '@mui/icons-material/Search';
import {useDialogs} from '@toolpad/core';
import SubscribeDialog from './SubscribeDialog';
import {MAIN_PAGE_IDS} from '../../constants/test_ids';

type FormData = {
    city: string;
};

type WeatherProps = Weather & {
    city: string;
    onCitySearch?: (city: string) => void;
};

const WeatherBlock: FC<WeatherProps> = ({temperature, humidity, description, city, onCitySearch}) => {
    const [subscribeOpen, setSubscribeOpen] = useState(false);
    const dialogs = useDialogs();
    const form = useForm<FormData>({
        reValidateMode: 'onChange',
        mode: 'all',
        defaultValues: {city: ''},
    });

    const {handleSubmit, control} = form;

    const onSubmit = (data: FormData) => {
        if (onCitySearch) {
            onCitySearch(data.city.trim());
        }
    };
    return (
        <Stack direction={'column'} spacing={3}>
            <form onSubmit={handleSubmit(onSubmit)}>
                <TextFieldElement
                    rules={{
                        required: 'Field required!',
                    }}
                    size="small"
                    type="text"
                    name="city"
                    label="Search city ..."
                    control={control}
                    fullWidth
                    slotProps={{
                        input: {
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        onClick={handleSubmit(onSubmit)}
                                        edge="end"
                                        size="small"
                                        data-testid={MAIN_PAGE_IDS.searchButton}
                                    >
                                        <SearchIcon />
                                    </IconButton>
                                </InputAdornment>
                            ),
                        },
                    }}
                    data-testid={MAIN_PAGE_IDS.searchInput}
                />
            </form>
            <Stack direction={'row'} spacing={2}>
                <Card variant="outlined" sx={{flex: 1}}>
                    <CardContent>
                        <Typography variant="caption">Temperature</Typography>
                        <Typography variant="h4" data-testid={MAIN_PAGE_IDS.tempValue}>
                            {temperature}°C
                        </Typography>
                    </CardContent>
                </Card>
                <Card variant="outlined" sx={{flex: 1}}>
                    <CardContent>
                        <Typography variant="caption">Humidity</Typography>
                        <Typography variant="h4" data-testid={MAIN_PAGE_IDS.humidityValue}>
                            {humidity}%
                        </Typography>
                    </CardContent>
                </Card>
            </Stack>
            <Card variant="outlined">
                <CardContent>
                    <Typography variant="caption">Description for weather in {city}</Typography>
                    <Typography variant="h4" data-testid={MAIN_PAGE_IDS.descriptionValue}>
                        {description}
                    </Typography>
                </CardContent>
            </Card>
            <Box display="flex" justifyContent="center" mt={2}>
                <Button
                    data-testid={MAIN_PAGE_IDS.subscribeButton}
                    variant="contained"
                    color="primary"
                    onClick={() => setSubscribeOpen(true)}
                >
                    Subscribe for updates
                </Button>
            </Box>
            <SubscribeDialog open={subscribeOpen} onClose={() => setSubscribeOpen(false)} />
        </Stack>
    );
};

export default memo(WeatherBlock);
