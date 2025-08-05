import {v4 as uuidv4} from 'uuid';

export function randAlfaNumeric(length: number = 15): string {
    const characters: string = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
    let result: string = '';

    for (let i: number = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * characters.length));
    }
    return result;
}

export function generateAutotestEmail(length: number = 10): string {
    return `autotest_${randAlfaNumeric(length)}@testmail.com`;
}

export function generateUUID(): string {
    return uuidv4();
}
