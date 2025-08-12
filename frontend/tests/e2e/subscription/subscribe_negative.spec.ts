import {test, expect, Page} from '@playwright/test';
import {MainPage} from '../../pom/mainPage';
import {User} from '../../core/models/user';
import {Subscription} from '../../core/models/subscription';
import {Context} from '../../fixtures/context';
import {SubscribeDialog} from '../../pom/subscribeDialog';

test.describe('Subscribe negative flow', () => {
    let ctx: Context;
    let mainPage: MainPage;
    let subDialog: SubscribeDialog;
    let existingUser: User;

    test.beforeEach(async ({page}) => {
        existingUser = new User({
            subscriptions: [new Subscription({city: 'Kyiv', isConfirmed: true})],
        });

        ctx = new Context({users: [existingUser]});
        await ctx.performSetup();

        mainPage = new MainPage(page);
        subDialog = mainPage.subDialog;

        await mainPage.loadMainPage();
    });

    test.afterEach(async () => {
        await ctx.performTeardown();
    });

    test('Should not allow double subscription', async ({page}) => {
        await mainPage.openSubscribeDialog();

        await expect(mainPage.subDialog.subscribeDialog).toBeVisible();
        await subDialog.subscribe({city: 'Kharkiv', email: existingUser.email});

        await expect(mainPage.toastAlert).toBeVisible();
        await expect(mainPage.toastAlert).toHaveText('Failed to subscribe');
        await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorError/);
    });
});
