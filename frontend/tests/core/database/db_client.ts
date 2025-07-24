import {Client} from 'pg';
import {User} from '../models/user';
import {Subscription} from '../models/subscription';

type DbOptions = {
    uri: string;
};

export class DatabaseClient {
    private db: Client;

    constructor(options: DbOptions) {
        this.db = new Client({connectionString: options.uri});
    }

    public async connect(): Promise<void> {
        await this.db.connect();
    }

    public async close(): Promise<void> {
        await this.db.end();
    }

    public async insertUsers(users: User[]): Promise<void> {
        if (users.length === 0) return;
        const values: unknown[] = [];
        const placeholders = users.map((u, i) => {
            const idx = i * 4;
            values.push(u.email, u.createdAt, u.updatedAt, u.deletedAt ?? null);
            return `($${idx + 1}, $${idx + 2}, $${idx + 3}, $${idx + 4})`;
        });

        const {rows} = await this.db.query(
            `INSERT INTO users (email, created_at, updated_at, deleted_at)
           VALUES ${placeholders.join(',')}
           RETURNING id`,
            values,
        );

        if (rows.length !== users.length) {
            throw new Error('Mismatch between inserted users and returned rows');
        }

        rows.forEach((row, i) => {
            const user = users[i];
            user.id = row.id;
            for (const sub of user.subscriptions) {
                sub.userId = user.id;
            }
        });
    }

    public async insertSubscriptions(subs: Subscription[]): Promise<void> {
        if (!subs.length) return;
        const values: unknown[] = [];
        const placeholders = subs.map((s, i) => {
            const idx = i * 10;
            values.push(
                s.city,
                s.frequency,
                s.userId,
                s.isConfirmed,
                s.confirmToken,
                s.tokenExpires,
                s.confirmedAt ?? null,
                s.createdAt,
                s.updatedAt,
                s.deletedAt ?? null,
            );
            return `($${idx + 1}, $${idx + 2}, $${idx + 3}, $${idx + 4}, $${idx + 5},
                   $${idx + 6}, $${idx + 7}, $${idx + 8}, $${idx + 9}, $${idx + 10})`;
        });

        const {rows} = await this.db.query(
            `INSERT INTO subscriptions
             (city, frequency, user_id, is_confirmed, confirm_token,
              token_expires, confirmed_at, created_at, updated_at, deleted_at)
           VALUES ${placeholders.join(',')}
           RETURNING id`,
            values,
        );

        rows.forEach((r, i) => {
            subs[i].id = r.id;
        });
    }

    public async cleanup(users: User[]): Promise<void> {
        if (users.length === 0) return;

        const emails = users.map(u => u.email).filter((e): e is string => !!e);

        if (emails.length === 0) return;

        await this.db.query(`DELETE FROM subscriptions WHERE user_id IN (SELECT id FROM users WHERE email = ANY($1))`, [emails]);

        await this.db.query(`DELETE FROM users WHERE email = ANY($1)`, [emails]);
    }
}
