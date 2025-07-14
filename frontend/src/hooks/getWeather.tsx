import {useQuery} from '@tanstack/react-query';
import {Weather, WeatherService} from '../api';
import {useNotifications} from "@toolpad/core";

export const getWeather = (city: string) => {
    const {
        data,
        isLoading,
        isPending,
        isRefetching,
        isFetching,
        isError,
        error,
    } = useQuery<Weather, Error>({
        queryKey: ['weather', city],
        queryFn: () => WeatherService.getWeather(city),
        enabled: !!city,
        refetchOnWindowFocus: false,
        retry: 0,
    });

    return {
        weather: data,
        loading: isLoading,
        pending: isPending || isRefetching || isFetching,
        isError,
        error

    };
};
