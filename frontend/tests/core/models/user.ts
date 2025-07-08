import {generateAutotestEmail} from '../utils/generate_utils';
import {Subscription} from './subscription';

export type UserOptions = {
    id?: number;
    email?: string;
    subscriptions?: Subscription[];
    createdAt?: Date;
    updatedAt?: Date;
    deletedAt?: Date;
};

export class User {
    public id?: number;
    public email: string;
    public subscriptions?: Subscription[];
    public createdAt: Date;
    public updatedAt: Date;
    public deletedAt?: Date;

    constructor(options: UserOptions) {
        this.id = options.id;
        this.email = options.email || generateAutotestEmail();
        this.subscriptions = options.subscriptions;
        this.createdAt = options.createdAt || new Date(Date.now());
        this.updatedAt = options.updatedAt || new Date(Date.now());
        this.deletedAt = options.deletedAt;
    }
}
