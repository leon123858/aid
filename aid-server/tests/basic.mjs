import {execSync, spawn} from 'child_process';
import crypto from 'crypto';
import assert from 'assert';

const serverUrl = 'http://127.0.0.1:8080';
// const serverUrl = 'http://20.2.209.109';
let serverProcess;
let aid;
let privateKey;
let publicKey;

const startServer = () => {
    serverProcess = spawn('../bin/aid', ['server']);
    execSync('sleep 2');
};

const stopServer = () => {
    serverProcess.kill();
};

const generateRSAKeyPair = () => {
    return crypto.generateKeyPairSync('rsa', {
        modulusLength: 2048,
        publicKeyEncoding: {type: 'spki', format: 'pem'},
        privateKeyEncoding: {type: 'pkcs8', format: 'pem'},
    });
};

const addUser = async () => {
    const {privateKey: privKey, publicKey: pubKey} = generateRSAKeyPair();
    aid = crypto.randomUUID();
    const response = await fetch(`${serverUrl}/api/register`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            aid: aid,
            publicKey: pubKey,
            ip: '127.0.1.3',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    privateKey = privKey;
    publicKey = pubKey;
};

async function testRegister() {
    console.log('Testing Register...');
    const {publicKey} = generateRSAKeyPair();
    const response = await fetch(`${serverUrl}/api/register`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            aid: crypto.randomUUID(),
            publicKey: publicKey,
            ip: '127.0.0.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
}

async function testRegisterWithInvalidPublicKey() {
    console.log('Testing Register with invalid public key...');
    const response = await fetch(`${serverUrl}/api/register`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            aid: crypto.randomUUID(),
            publicKey: 'invalid public key',
            ip: '127.0.0.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 400);
    assert.deepStrictEqual(data, {
        result: false,
        content: 'invalid public key',
    });
}

async function testLogin() {
    console.log('Testing Login...');
    await addUser();
    const timestamp = Date.now().toString();
    const sign = crypto.sign('sha256', Buffer.from(timestamp), {
        key: privateKey,
        padding: crypto.constants.RSA_PKCS1_PADDING,
    });
    const response = await fetch(`${serverUrl}/api/login`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            aid: aid,
            sign: sign.toString('base64'),
            timestamp: timestamp,
            ip: '127.0.1.3',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
    return data.content; // Return token for Ask test
}

async function testAsk(token) {
    console.log('Testing Ask...');
    const response = await fetch(`${serverUrl}/api/ask`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${token}`
        },
        body: JSON.stringify({
            ip: '127.0.1.3',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
    return data.content; // Return uid for Check and Verify tests
}

async function testCheckOnline(uid) {
    console.log('Testing Check - Online...');
    const response = await fetch(`${serverUrl}/api/check`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            uid: uid,
            ip: '127.0.1.3',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.strictEqual(data.content, 'online');
}

async function testCheckOffline(uid) {
    console.log('Testing Check - Offline...');
    const response = await fetch(`${serverUrl}/api/check`, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({
            uid: uid,
            ip: '127.0.1.3',
            browser: 'Safari', // Different browser should result in offline status
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.strictEqual(data.content, 'offline');
}

async function testVerifySuccess(token, uid) {
    console.log('Testing Verify - Success...');
    const response = await fetch(`${serverUrl}/api/verify`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': token
        },
        body: JSON.stringify({
            uid: uid,
            ip: '127.0.0.1', // Different IP is allowed for verify
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
}

async function testVerifyInvalidToken(uid) {
    console.log('Testing Verify - Invalid Token...');
    const response = await fetch(`${serverUrl}/api/verify`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'invalid_token'
        },
        body: JSON.stringify({
            uid: uid,
            ip: '127.0.0.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 401);
    assert.strictEqual(data.result, false);
    assert.strictEqual(data.content, 'invalid token');
}

async function runTests() {
    try {
        startServer();
        await testRegister();
        await testRegisterWithInvalidPublicKey();
        const token = await testLogin();
        const uid = await testAsk(token);
        await testCheckOnline(uid);
        await testCheckOffline(uid);
        await testVerifySuccess(token, uid);
        await testVerifyInvalidToken(uid);
        console.log('All tests passed successfully!');
    } catch (error) {
        console.error('Test failed:', error);
    } finally {
        stopServer();
    }
}

runTests().then(r =>
    console.log('Tests completed')
).catch(e =>
    console.error('Tests failed:', e)
);