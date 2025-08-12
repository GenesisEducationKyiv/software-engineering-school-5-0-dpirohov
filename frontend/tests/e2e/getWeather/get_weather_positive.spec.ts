import {test, expect} from '@playwright/test';
import {MainPage} from '../../pom/mainPage';

test.describe('Get weather positive flow', () => {
    let mainPage: MainPage;

    test.beforeEach(async ({page}) => {
        mainPage = new MainPage(page);

        await mainPage.loadMainPage();
    });

    test('Should load and display weather data for searched city', async ({page}) => {
        const mockWeatherResponse = {
            temperature: 25,
            humidity: 60,
            description: 'Sunny',
            city: 'Poltava',
        };

        await page.route('**/api/v1/weather?city=Poltava', route =>
            route.fulfill({
                status: 200,
                contentType: 'application/json',
                body: JSON.stringify(mockWeatherResponse),
            }),
        );

        await mainPage.searchInput.fill('Poltava');
        await mainPage.searchButton.click();

        await expect(mainPage.tempValue).toHaveText(`${mockWeatherResponse.temperature}°C`);
        await expect(mainPage.humidityValue).toHaveText(`${mockWeatherResponse.humidity}%`);
        await expect(mainPage.descriptionValue).toHaveText(mockWeatherResponse.description);
    });
});
