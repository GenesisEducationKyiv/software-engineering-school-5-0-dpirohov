import {test, expect} from '@playwright/test';

import {MainPage} from '../../pom/mainPage';
import {Context} from '../../fixtures/context';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';
import {ConfirmationPage} from '../../pom/confirmationPage';
import {generateUUID} from '../../core/utils/generate_utils';

test.describe('Confirm subscription negative flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let confirmationPage: ConfirmationPage;
    let existingUser: User;

    test.beforeAll(async () => {
        existingUser = new User({
            subscriptions: [new Subscription({city: 'Kyiv', isConfirmed: true})],
        });

        ctx = new Context({users: [existingUser]});
        await ctx.performSetup();
    });

    test.beforeEach(async ({page}) => {
        mainPage = new MainPage(page);
        confirmationPage = new ConfirmationPage(page);
    });

    test.afterAll(async () => {
        await ctx.performTeardown();
    });

    const tokens = [
        {name: 'already used token', token: () => existingUser.subscriptions[0].confirmToken},
        {name: 'invalid token', token: generateUUID},
    ];

    for (const {name, token} of tokens) {
        test(`Should not confirm with ${name}`, async ({page}) => {
            await page.goto(`/confirm/${token()}`);

            await expect(confirmationPage.confirmationText).toBeVisible();
            await expect(confirmationPage.linkToMainPage).toBeVisible();
            await expect(confirmationPage.confirmationText).toHaveText(/Confirmation failed ❌/);

            await expect(mainPage.toastAlert).toBeVisible();
            await expect(mainPage.toastAlert).toHaveText(/Confirmation failed/i);
            await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorError/);

            await confirmationPage.linkToMainPage.click();
            await expect(page.getByText('Kyiv')).toBeVisible();
            await expect(mainPage.subscribeButton).toBeVisible();
        });
    }
});
