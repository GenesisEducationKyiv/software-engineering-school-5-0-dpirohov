import {Box} from '@mui/material';

import {getWeather} from '../hooks/getWeather';
import LoadingSkeleton from './components/LoadingSkeleton';
import WeatherBlock from './components/Weather';
import {useEffect, useState} from "react";
import {useNotifications} from "@toolpad/core";


export default function WeatherPage() {
    const [city, setCity] = useState('Kyiv');
    const {weather, loading, pending, isError, error} = getWeather(city);
    const notifications = useNotifications();

    useEffect(() => {
        if (isError && error) {
            notifications.show(`Failed to load weather: ${error.message}`, {
                severity: 'error',
            });
        }
    }, [isError, error]);

    return (
        <Box display="flex" justifyContent="center" alignItems="center">
            <Box maxWidth={800} width="100%" px={2}>
                {loading || pending ? (
                    <LoadingSkeleton />
                ) : (
                    <WeatherBlock
                        temperature={weather?.temperature ?? 0}
                        humidity={weather?.humidity ?? 0}
                        description={weather?.description ?? "No description provided :("}
                        city={city}
                        onCitySearch={setCity}
                    />
                )}
            </Box>
        </Box>
    );
}
