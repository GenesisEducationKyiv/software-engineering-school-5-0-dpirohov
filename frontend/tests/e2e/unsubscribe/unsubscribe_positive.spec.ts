import {test, expect} from '@playwright/test';

import {MainPage} from '../../pom/mainPage';
import {Context} from '../../fixtures/context';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';
import {UnsubscribePage} from '../../pom/unsubscribePage';

test.describe('Unsubscribe from updates, positive flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let unsubscribePage: UnsubscribePage;
    let subscription: Subscription;

    test.beforeAll(async () => {
        subscription = new Subscription({city: 'Kyiv', isConfirmed: true});
        const existingUser = new User({
            subscriptions: [subscription],
        });

        ctx = new Context({users: [existingUser]});
        await ctx.performSetup();
    });

    test.beforeEach(async ({page}) => {
        mainPage = new MainPage(page);
        unsubscribePage = new UnsubscribePage(page);
    });

    test.afterAll(async () => {
        await ctx.performTeardown();
    });

    test(`Should unsubscribe with valid confirmed token`, async ({page}) => {
        await page.goto(`/unsubscribe/${subscription.confirmToken}`);

        await expect(unsubscribePage.confirmationText).toBeVisible();
        await expect(unsubscribePage.linkToMainPage).toBeVisible();
        await expect(unsubscribePage.confirmationText).toHaveText(/Succsessfuly unsubscribed ✅/);

        await expect(mainPage.toastAlert).toBeVisible();
        await expect(mainPage.toastAlert).toHaveText(/Succsessfuly unsubscribed!/i);
        await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorSuccess/);

        await unsubscribePage.linkToMainPage.click();
        await expect(page.getByText('Kyiv')).toBeVisible();
        await expect(mainPage.subscribeButton).toBeVisible();
    });
});
