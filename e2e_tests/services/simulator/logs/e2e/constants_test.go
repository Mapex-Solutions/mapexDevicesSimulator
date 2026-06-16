package e2e

// listFixtureCount is the number of log entries the cursor-walk test seeds. 15
// at limit=1 exposes the pagination edges (first page, mid pages, last page)
// that a smaller count would not.
const listFixtureCount = 15

// nonExistentDeviceID is a deviceId no fixture uses, for the filter no-match
// tests.
const nonExistentDeviceID = "no-such-device-zzzz"
