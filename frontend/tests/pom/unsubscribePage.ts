import {BasePage} from './basePage';
import {UNSUBSCRIBE_PAGE_IDS} from '../../src/constants/test_ids';

export class UnsubscribePage extends BasePage {
    public get confirmationText() {
        return this.page.getByTestId(UNSUBSCRIBE_PAGE_IDS.confirmation);
    }

    public get linkToMainPage() {
        return this.page.getByTestId(UNSUBSCRIBE_PAGE_IDS.linkToMainPage);
    }
}
