import * as dotenv from 'dotenv';
import * as path from 'path';

export class Config {
    public readonly dbUrl: string;
    public readonly baseUrl: string;

    constructor() {
        this.loadEnv();
        this.dbUrl = this.requireEnv('DB_URL');
        this.baseUrl = this.requireEnv('BASE_URL');
    }

    private requireEnv(key: string): string {
        const value = process.env[key];
        if (!value) {
            throw new Error(`Missing required environment variable: ${key}`);
        }
        return value;
    }

    private loadEnv(): void {
        const envPath = path.resolve(__dirname, '../tests/.env');
        dotenv.config({path: envPath});
    }
}

export const CONFIG = new Config();
