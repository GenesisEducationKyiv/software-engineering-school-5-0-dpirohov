import { test, expect } from '@playwright/test';
import { MainPage } from '../../pom/mainPage';

type ErrorScenario = {
  status: number;
  responseBody: Record<string, any>;
  errorMessage: string;
  errorTextInToast: string;
};

const errorScenarios: ErrorScenario[] = [
  {
    status: 404,
    responseBody: {
      statusCode: 404,
      message: 'City not found',
      error: 'Not Found',
    },
    errorMessage: 'Failed to load weather: City not found',
    errorTextInToast: 'City not found',
  },
  {
    status: 400,
    responseBody: {
      statusCode: 400,
      message: 'Invalid Request',
      error: 'Bad Request',
    },
    errorMessage: 'Failed to load weather: Invalid request',
    errorTextInToast: 'Invalid request',
  },
  {
    status: 500,
    responseBody: {
      statusCode: 500,
      message: 'Internal server error',
      error: 'Internal Server Error',
    },
    errorMessage: 'Failed to load weather: Internal Server Error',
    errorTextInToast: 'Internal Server Error',
  },
];

test.describe('Get weather negative flow', () => {
  let mainPage: MainPage;

  test.beforeEach(async ({ page }) => {
    mainPage = new MainPage(page);
    await mainPage.loadMainPage();
  });

  for (const scenario of errorScenarios) {
    test(`Mock error ${scenario.status} - "${scenario.responseBody.message}"`, async ({ page }) => {
      await page.route(`**/api/v1/weather?city=UnknownCity`, route => {
        route.fulfill({
          status: scenario.status,
          contentType: 'application/json',
          body: JSON.stringify(scenario.responseBody),
        });
      });

      await mainPage.searchInput.fill('UnknownCity');
      await mainPage.searchButton.click();

      await expect(mainPage.toastAlert).toBeVisible();
      await expect(mainPage.toastAlert).toHaveText(scenario.errorMessage);
      await expect(mainPage.toastAlert).toHaveClass(/MuiAlert-colorError/);
    });
  }
});
