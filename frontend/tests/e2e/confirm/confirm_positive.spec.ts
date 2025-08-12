import {test, expect} from '@playwright/test';

import {MainPage} from '../../pom/mainPage';
import {Context} from '../../fixtures/context';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';
import {ConfirmationPage} from '../../pom/confirmationPage';

test.describe('Confirm subscription positive flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let confirmationPage: ConfirmationPage;
    let existingUser: User;

    test.beforeAll(async () => {
        existingUser = new User({
            subscriptions: [new Subscription({city: 'Kyiv'})],
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

    test('Should confirm subscription token', async ({page}) => {
        await page.goto(`/confirm/${existingUser.subscriptions[0].confirmToken}`);
        await expect(confirmationPage.confirmationText).toBeVisible();
        await expect(confirmationPage.linkToMainPage).toBeVisible();
        await expect(confirmationPage.confirmationText).toHaveText('Subscription confirmed ✅');

        await expect(mainPage.toastAlert).toBeVisible();
        await expect(mainPage.toastAlert).toHaveText('Subscription confirmed!');
        await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorSuccess/);

        // confirm link is valid
        await confirmationPage.linkToMainPage.click();
        await expect(page.getByText('Kyiv')).toBeVisible();
        await expect(mainPage.subscribeButton).toBeVisible();
    });
});
