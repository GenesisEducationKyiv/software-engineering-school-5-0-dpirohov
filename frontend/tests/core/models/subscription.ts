import {generateUUID} from '../utils/generate_utils';

const minute = 60 * 1000;

export type Frequency = 'daily' | 'hourly' | 'weekly';

export type SubscriptionOptions = {
    id?: number;
    city: string;
    frequency?: Frequency;
    userId?: number;
    isConfirmed?: boolean;
    confirmToken?: string;
    tokenExpires?: Date;
    confirmedAt?: Date;
    createdAt?: Date;
    updatedAt?: Date;
    deletedAt?: Date;
};

export class Subscription {
    public id?: number;
    public city: string;
    public frequency: Frequency;
    public userId: number;
    public isConfirmed: boolean;
    public confirmToken: string;
    public tokenExpires: Date;
    public confirmedAt?: Date;
    public createdAt: Date;
    public updatedAt: Date;
    public deletedAt?: Date;

    constructor(options: SubscriptionOptions) {
        this.id = options.id;
        this.city = options.city;
        this.frequency = options.frequency || 'daily';
        this.userId = options.userId;
        this.isConfirmed = options.isConfirmed || false;
        this.confirmToken = options.confirmToken || generateUUID();
        this.tokenExpires = options.tokenExpires || new Date(Date.now() + 15 * minute);
        this.createdAt = options.createdAt || new Date(Date.now());
        this.updatedAt = options.updatedAt || new Date(Date.now());
        this.deletedAt = options.deletedAt;
    }
}
