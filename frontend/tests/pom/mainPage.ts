import {expect, Page} from '@playwright/test';
import {MAIN_PAGE_IDS, SUBSCRIBE_DIALOG_IDS} from '../../src/constants/test_ids';
import {BasePage} from './basePage';
import {SubscribeDialog} from './subscribeDialog';

export class MainPage extends BasePage {
    public subDialog: SubscribeDialog;

    constructor(page: Page) {
        super(page);
        this.subDialog = new SubscribeDialog(page);
    }

    get searchInput() {
        return this.page.locator(`[data-testid="${MAIN_PAGE_IDS.searchInput}"] input`);
    }
    get searchButton() {
        return this.page.getByTestId(MAIN_PAGE_IDS.searchButton);
    }
    get tempValue() {
        return this.page.getByTestId(MAIN_PAGE_IDS.tempValue);
    }
    get humidityValue() {
        return this.page.getByTestId(MAIN_PAGE_IDS.humidityValue);
    }
    get descriptionValue() {
        return this.page.getByTestId(MAIN_PAGE_IDS.descriptionValue);
    }
    get subscribeButton() {
        return this.page.getByTestId(MAIN_PAGE_IDS.subscribeButton);
    }

    async loadMainPage() {
        await this.page.goto('/');
        await expect(this.page).toHaveTitle('Weather Api service');
        await expect(this.page.getByText('Kyiv')).toBeVisible();
        await expect(this.subscribeButton).toBeVisible();
    }

    async openSubscribeDialog() {
        await this.subscribeButton.click();
        await this.subDialog.subscribeDialog.waitFor({state: 'visible'});
    }
}
