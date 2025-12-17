// DEVELOPMENT ONLY - Mock initData for testing outside Telegram
// DO NOT USE IN PRODUCTION!

// Uncomment this function and import it in logs.ts to test without Telegram
export function getMockInitData(): string {
  // Mock data for user ID 12345
  // Format: URL-encoded JSON: {"id":12345,"first_name":"Test","username":"testuser"}
  return 'user=%7B%22id%22%3A12345%2C%22first_name%22%3A%22Test%22%2C%22username%22%3A%22testuser%22%7D';
}

// Instructions to use:
// 1. Open web/src/api/logs.ts
// 2. Replace getInitData() function with:
//
//    import { getMockInitData } from './logs-dev';
//
//    function getInitData(): string {
//      return getMockInitData();  // DEV ONLY
//    }
//
// 3. Reload the app
//
// REMEMBER TO REVERT THIS BEFORE PRODUCTION!
