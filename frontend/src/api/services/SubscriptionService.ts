/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
import type {CancelablePromise} from '../core/CancelablePromise';
import {OpenAPI} from '../core/OpenAPI';
import {request as __request} from '../core/request';
export class SubscriptionService {
    /**
     * Subscribe to weather updates
     * Subscribe an email to receive weather updates for a specific city with chosen frequency.
     * @param email Email address to subscribe
     * @param city City for weather updates
     * @param frequency Frequency of updates (hourly or daily)
     * @returns any Subscription successful. Confirmation email sent.
     * @throws ApiError
     */
    public static subscribe(email: string, city: string, frequency: 'hourly' | 'daily'): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'POST',
            url: 'api/v1/subscribe',
            mediaType: 'application/json',
            body: {
                email: email,
                city: city,
                frequency: frequency,
            },
            errors: {
                400: `Invalid input`,
                409: `Email already subscribed`,
            },
        });
    }
    /**
     * Confirm email subscription
     * Confirms a subscription using the token sent in the confirmation email.
     * @param token Confirmation token
     * @returns any Subscription confirmed successfully
     * @throws ApiError
     */
    public static confirmSubscription(token: string): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: 'api/v1/confirm/{token}',
            path: {
                token: token,
            },
            errors: {
                400: `Invalid token`,
                404: `Token not found`,
            },
        });
    }
    /**
     * Unsubscribe from weather updates
     * Unsubscribes an email from weather updates using the token sent in emails.
     * @param token Unsubscribe token
     * @returns any Unsubscribed successfully
     * @throws ApiError
     */
    public static unsubscribe(token: string): CancelablePromise<any> {
        return __request(OpenAPI, {
            method: 'GET',
            url: 'api/v1/unsubscribe/{token}',
            path: {
                token: token,
            },
            errors: {
                400: `Invalid token`,
                404: `Token not found`,
            },
        });
    }
}
