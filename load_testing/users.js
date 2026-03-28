import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { SharedArray } from 'k6/data';
import { Counter, Rate, Trend } from 'k6/metrics';

const errorRate      = new Rate('errors');
const createDuration = new Trend('duration_create_user', true);
const getDuration    = new Trend('duration_get_user',    true);
const updateDuration = new Trend('duration_update_user', true);
const deleteDuration = new Trend('duration_delete_user', true);
const notFoundCount  = new Counter('user_not_found');

export const options = {
    scenarios: {
        // Scenario 1 — Casual browsers: register, peek, leave
        casual_browsers: {
            executor:        'ramping-vus',
            startVUs:        0,
            stages: [
                { duration: '1m',  target: 10 },
                { duration: '3m',  target: 10 },
                { duration: '30s', target: 0  },
            ],
            exec: 'casualBrowser',
            tags: { scenario: 'casual' },
        },

        // Scenario 2 — Power users: full CRUD + repeat reads
        power_users: {
            executor:        'ramping-vus',
            startVUs:        0,
            stages: [
                { duration: '1m',  target: 5  },
                { duration: '3m',  target: 5  },
                { duration: '30s', target: 0  },
            ],
            exec: 'powerUser',
            tags: { scenario: 'power' },
        },

        // Scenario 3 — Read-heavy traffic (already-existing users)
        read_traffic: {
            executor:        'constant-arrival-rate',
            rate:            20,
            timeUnit:        '1s',
            duration:        '4m',
            preAllocatedVUs: 10,
            maxVUs:          30,
            exec: 'readOnlyUser',
            tags: { scenario: 'read_heavy' },
        },
    },

    thresholds: {
        http_req_failed:      ['rate<0.02'],
        http_req_duration:    ['p(95)<500'],
        duration_create_user: ['p(95)<600'],
        duration_get_user:    ['p(95)<300'],
        errors:               ['rate<0.05'],
    },
};

// Shared state pool (pre-created user IDs for read scenario)
// K6 SharedArray is read-once at init time, so we seed a small static pool.
// In real use, replace with a CSV / JSON fixture generated beforehand.
const seedUsers = new SharedArray('seed_users', function () {
    // Generates deterministic usernames the read scenario can attempt to look up.
    // These won't exist unless the setup() creates them — see setup() below.
        return Array.from({ length: 50 }, (_, i) => ({ id: `seed_${i}` }));
});


const BASE_URL = __ENV.BASE_URL || 'http://arch.homework:8080';
const JSON_HEADERS = { headers: { 'Content-Type': 'application/json' } };

function uniqueTag() {
    return `${__VU}_${__ITER}_${Date.now()}`;
}

function buildUserPayload(tag, prefix = '') {
    return JSON.stringify({
        username:  `${prefix}user_${tag}`,
        firstName: `First_${tag}`,
        lastName:  `Last_${tag}`,
        email:     `${prefix}${tag}@example.com`,
        phone:     `+7900${String(Math.floor(Math.random() * 9999999)).padStart(7, '0')}`,
    });
}

/** Realistic think-time: short pause drawn from a skewed distribution. */
function thinkTime(minSec = 0.3, maxSec = 2.5) {
    sleep(minSec + Math.random() * (maxSec - minSec));
}

/**
 * Creates a user and returns { userId, username } or null on failure.
 */
function createUser(tag, prefix = '') {
    const res = http.post(`${BASE_URL}/user`, buildUserPayload(tag, prefix), JSON_HEADERS);
    createDuration.add(res.timings.duration);

    const ok = check(res, {
        'create: 200': (r) => r.status === 200,
        'create: has id': (r) => {
            try { return r.json('id') !== undefined; } catch { return false; }
        },
    });

    errorRate.add(!ok);
    if (!ok) return null;

    return { userId: res.json('id'), username: `${prefix}user_${tag}` };
}

/**
 * GETs a user by id; returns parsed body or null.
 */
function getUser(userId) {
    const res = http.get(`${BASE_URL}/user/${userId}`);
    getDuration.add(res.timings.duration);

    if (res.status === 404) { notFoundCount.add(1); return null; }

    const ok = check(res, { 'get: 200': (r) => r.status === 200 });
    errorRate.add(!ok);
    return ok ? res.json() : null;
}

/**
 * PUTs updated profile for a user.
 */
function updateUser(userId, tag) {
    const payload = JSON.stringify({
        username:  `user_${tag}`,
        firstName: `Updated_${tag}`,
        lastName:  `Updated_${tag}`,
        email:     `${tag}_v2@example.com`,
        phone:     `+7800${String(Math.floor(Math.random() * 9999999)).padStart(7, '0')}`,
    });

    const res = http.put(`${BASE_URL}/user/${userId}`, payload, JSON_HEADERS);
    updateDuration.add(res.timings.duration);

    const ok = check(res, { 'update: 200': (r) => r.status === 200 });
    errorRate.add(!ok);
    return ok;
}

/**
 * DELETE user.
 */
function deleteUser(userId) {
    const res = http.del(`${BASE_URL}/user/${userId}`);
    deleteDuration.add(res.timings.duration);

    const ok = check(res, { 'delete: 204': (r) => r.status === 204 });
    errorRate.add(!ok);
    return ok;
}

// Setup: create a pool of persistent users for the read scenario
export function setup() {
    const persistentUsers = [];
    for (let i = 0; i < 20; i++) {
        const tag = `seed_${i}_${Date.now()}`;
        const user = createUser(tag, 'persistent_');
        if (user) persistentUsers.push(user);
        sleep(0.05); // avoid burst on setup
    }
    console.log(`Setup created ${persistentUsers.length} persistent users.`);
    return { persistentUsers };
}

// Teardown: clean up persistent users
export function teardown(data) {
    for (const user of data.persistentUsers) {
        deleteUser(user.userId);
    }
    console.log('Teardown: persistent users deleted.');
}

// Scenario A: Casual Browser
// Registers, quickly checks their profile once, then disappears.
export function casualBrowser() {
    const tag = uniqueTag();

    group('casual: register', () => {
        const user = createUser(tag, 'casual_');
        if (!user) return;

        thinkTime(1, 3);
        group('casual: peek at profile', () => {
            getUser(user.userId);
        });

        thinkTime(2, 6);
        // 30 % of casual users immediately delete their account
        if (Math.random() < 0.3) {
            group('casual: rage-quit', () => {
                deleteUser(user.userId);
            });
        }
        // The rest leave their account alive (realistic churn)
    });
}

// Scenario B: Power User
// Full lifecycle: create -> read x N -> update -> read again -> delete
export function powerUser() {
    const tag = uniqueTag();

    group('power: full lifecycle', () => {
        // 1. Register
        const user = createUser(tag, 'power_');
        if (!user) return;
        thinkTime(0.5, 1.5);

        // 2. Read their profile 2-4 times (tabs, refreshes)
        const reads = 2 + Math.floor(Math.random() * 3);
        for (let i = 0; i < reads; i++) {
            group(`power: read #${i + 1}`, () => getUser(user.userId));
            thinkTime(0.3, 1);
        }

        // 3. Update profile
        group('power: update', () => updateUser(user.userId, tag));
        thinkTime(0.5, 2);

        // 4. Verify update
        group('power: verify update', () => {
            const body = getUser(user.userId);
            check(body, {
                'power: updated firstName present': (b) =>
                    b && b.firstName && b.firstName.startsWith('Updated_'),
            });
        });
        thinkTime(1, 3);

        // 5. Clean up
        group('power: delete', () => deleteUser(user.userId));
    });
}

// Scenario C: Read-Only Traffic
// Simulates external services / other microservices reading user profiles.
// Uses the pool of persistent users created in setup().
export function readOnlyUser(data) {
    if (!data || !data.persistentUsers || data.persistentUsers.length === 0) {
        // Graceful degradation when setup pool is empty
        console.warn('readOnlyUser: no persistent users available');
        return;
    }

    const pool = data.persistentUsers;
    const target = pool[Math.floor(Math.random() * pool.length)];

    group('read: fetch user profile', () => {
        getUser(target.userId);
    });

    thinkTime(0.1, 0.5); // short pause - this is a service-to-service call
}