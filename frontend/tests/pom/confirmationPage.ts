import {BasePage} from './basePage';
import {CONFIRM_PAGE_IDS} from '../../src/constants/test_ids';

export class ConfirmationPage extends BasePage {
    public get confirmationText() {
        return this.page.getByTestId(CONFIRM_PAGE_IDS.confirmation);
    }

    public get linkToMainPage() {
        return this.page.getByTestId(CONFIRM_PAGE_IDS.linkToMainPage);
    }
}
