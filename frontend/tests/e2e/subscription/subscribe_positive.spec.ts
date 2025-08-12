import {test, expect} from '@playwright/test';

import {MainPage} from '../../pom/mainPage';
import {SubscribeDialog} from '../../pom/subscribeDialog';
import {Context} from '../../fixtures/context';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';

test.describe('Subscribe positive flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let subDialog: SubscribeDialog;
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
        subDialog = mainPage.subDialog;

        await mainPage.loadMainPage();
    });

    test.afterAll(async () => {
        await ctx.performTeardown();
    });

    test('Should subscribe new user', async () => {
        await mainPage.openSubscribeDialog();

        await expect(subDialog.subscribeDialog).toBeVisible();

        const newUser = ctx.createNewUser();

        await subDialog.subscribe({city: 'Kyiv', email: newUser.email});

        await expect(mainPage.toastAlert).toBeVisible();
        await expect(mainPage.toastAlert).toHaveText('Subscription submitted! Check your email for confirmation link!');
        await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorSuccess/);
    });

    test('Should subscribe existing user with not confirmed subscription', async () => {
        await mainPage.openSubscribeDialog();

        await expect(subDialog.subscribeDialog).toBeVisible();
        await subDialog.subscribe({city: 'Kharkiv', email: existingUser.email});

        await expect(mainPage.toastAlert).toBeVisible();
        await expect(mainPage.toastAlert).toHaveText('Subscription submitted! Check your email for confirmation link!');
        await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorSuccess/);
    });
});
