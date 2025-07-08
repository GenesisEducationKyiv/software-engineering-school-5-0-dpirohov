import {User} from '../core/models/user';
import {Subscription} from '../core/models/subscription';
import {CONFIG} from '../config';
import {DatabaseClient} from '../core/database/db_client';
import {generateAutotestEmail} from '../core/utils/generate_utils';

export type ContextOptions = {
    users?: User[];
};

export class Context {
    private db: DatabaseClient;
    public users: User[];
    private cleanUpUsers: User[];

    constructor(options: ContextOptions) {
        this.db = new DatabaseClient({uri: CONFIG.dbUrl});
        this.users = options.users || [];
        this.cleanUpUsers = [];
    }

    public async performTeardown(): Promise<void> {
        await this.cleanup();
        await this.db.close();
    }

    public async performSetup(): Promise<void> {
        await this.db.connect();
        await this.db.insertUsers(this.users);
        const subscriptions: Subscription[] = this.users.reduce((acc: Subscription[], user: User) => {
            acc.push(...(user.subscriptions ?? []));
            return acc;
        }, []);
        await this.db.insertSubscriptions(subscriptions);
    }

    public async cleanup(): Promise<void> {
        const allUsers = [...this.users, ...this.cleanUpUsers];
        await this.db.cleanup(allUsers);
    }

    public createNewUser(email?: string): User {
        if (!email) {
            email = generateAutotestEmail();
        }
        const user = new User({email});
        this.cleanUpUsers.push(user);
        return user;
    }
}
