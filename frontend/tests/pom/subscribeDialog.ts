import {BasePage} from './basePage';
import {SUBSCRIBE_DIALOG_IDS} from '../../src/constants/test_ids';

export type subscribeOptions = {
    city: string;
    email: string;
    frequency?: string;
};

export class SubscribeDialog extends BasePage {
    get subscribeDialog() {
        return this.page.getByTestId(SUBSCRIBE_DIALOG_IDS.dialog);
    }
    get subscribeCityInput() {
        return this.page.locator(`[data-testid="${SUBSCRIBE_DIALOG_IDS.cityInput}"] input`);
    }
    get subscribeEmailInput() {
        return this.page.locator(`[data-testid="${SUBSCRIBE_DIALOG_IDS.emailInput}"] input`);
    }
    get subscribeFrequencySelect() {
        return this.page.getByTestId(SUBSCRIBE_DIALOG_IDS.frequencySelect);
    }
    get subscribeDialogSubscribeButton() {
        return this.page.getByTestId(SUBSCRIBE_DIALOG_IDS.subscribeButton);
    }
    async subscribe(options: subscribeOptions) {
        await this.subscribeCityInput.fill(options.city);
        await this.subscribeEmailInput.fill(options.email);
        await this.selectMuiOption(this.subscribeFrequencySelect, options.frequency || 'Every day');
        await this.subscribeDialogSubscribeButton.click();
    }
}
