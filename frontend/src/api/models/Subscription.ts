/* generated using openapi-typescript-codegen -- do not edit */
/* istanbul ignore file */
/* tslint:disable */
/* eslint-disable */
export type Subscription = {
    /**
     * Email address
     */
    email: string;
    /**
     * City for weather updates
     */
    city: string;
    /**
     * Frequency of updates
     */
    frequency: Subscription.frequency;
    /**
     * Whether the subscription is confirmed
     */
    confirmed?: boolean;
};
export namespace Subscription {
    /**
     * Frequency of updates
     */
    export enum frequency {
        HOURLY = 'hourly',
        DAILY = 'daily',
    }
}
