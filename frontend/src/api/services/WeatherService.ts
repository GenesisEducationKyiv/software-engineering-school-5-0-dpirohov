/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {CancelablePromise} from '../core/CancelablePromise';
import {OpenAPI} from '../core/OpenAPI';
import {request as __request} from '../core/request';
export class WeatherService {
    /**
     * Get current weather for a city
     * Returns the current weather forecast for the specified city using WeatherAPI.com.
     * @param city City name for weather forecast
     * @returns any Successful operation - current weather forecast returned
     * @throws ApiError
     */
    public static getWeather(city: string): CancelablePromise<{
        /**
         * Current temperature
         */
        temperature?: number;
        /**
         * Current humidity percentage
         */
        humidity?: number;
        /**
         * Weather description
         */
        description?: string;
    }> {
        return __request(OpenAPI, {
            method: 'GET',
            url: 'api/v1/weather',
            query: {
                city: city,
            },
            errors: {
                400: `Invalid request`,
                404: `City not found`,
            },
        });
    }
}
