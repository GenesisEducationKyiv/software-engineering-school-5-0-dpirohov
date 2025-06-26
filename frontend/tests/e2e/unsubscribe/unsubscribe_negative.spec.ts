import {test, expect} from '@playwright/test';

import {MainPage} from '../../pom/mainPage';
import {Context} from '../../fixtures/context';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';
import {UnsubscribePage} from '../../pom/unsubscribePage';

test.describe('Confirm subscription negative flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let unsubscribePage: UnsubscribePage;
    let notConfirmedSubscription: Subscription;
    let deletedSubscription: Subscription;
    let notExistentSubscription: Subscription;

    test.beforeAll(async () => {
        notConfirmedSubscription = new Subscription({city: 'Kyiv', isConfirmed: false});
        deletedSubscription = new Subscription({city: 'Kyiv', isConfirmed: true, deletedAt: new Date(Date.now())});
        notExistentSubscription = new Subscription({city: 'Kyiv', isConfirmed: false});

        const existingUser = new User({
            subscriptions: [notConfirmedSubscription, deletedSubscription],
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

    const testData = [
        {
            description: 'Subscription not confirmed',
            token: () => notConfirmedSubscription.confirmToken,
            expectedText: 'Unsubscribe failed: Token not found',
        },
        {
            description: 'Subscription deleted',
            token: () => deletedSubscription.confirmToken,
            expectedText: 'Unsubscribe failed: Token not found',
        },
        {
            description: 'Token does not exist',
            token: () => notExistentSubscription.confirmToken,
            expectedText: 'Unsubscribe failed: Token not found',
        },
    ];

    for (const {description, token, expectedText} of testData) {
        test(`Should show corrent error when ${description}`, async ({page}) => {
            await page.goto(`/unsubscribe/${token}`);

            await expect(unsubscribePage.confirmationText).toBeVisible();
            await expect(unsubscribePage.linkToMainPage).toBeVisible();
            await expect(unsubscribePage.confirmationText).toHaveText(/Unsubscribe failed! ❌/);

            await expect(mainPage.toastAlert).toBeVisible();
            await expect(mainPage.toastAlert).toHaveText(expectedText);
            await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorError/);

            await unsubscribePage.linkToMainPage.click();
            await expect(page.getByText('Kyiv')).toBeVisible();
            await expect(mainPage.subscribeButton).toBeVisible();
        });
    }
});
