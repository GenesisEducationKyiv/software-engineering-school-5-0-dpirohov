import {Locator, Page} from '@playwright/test';

export class BasePage {
    readonly page: Page;

    constructor(page: Page) {
        this.page = page;
    }
    public async selectMuiOption(selectLocator: Locator, optionValue: string) {
        await selectLocator.click();
        const option = this.page.locator(`li[role="option"] >> text=${optionValue}`);
        await option.click();
    }

    public get toastAlert() {
        return this.page.getByRole('alert');
    }
}
