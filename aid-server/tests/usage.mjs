import {execSync, spawn} from 'child_process';
import assert from 'assert';
import crypto from 'crypto';

const serverUrl = 'http://127.0.0.1:8080';
let serverProcess;

const startServer = () => {
    serverProcess = spawn('../bin/aid', ['server']);
    execSync('sleep 2');
};

const stopServer = () => {
    if (serverProcess) {
        serverProcess.kill();
    }
};

let randomUUID = () => {
    return crypto.randomUUID();
};
let userName = randomUUID()
let publicKey = ""
let privateKey = ""
let aid = ""

const generateRSAKeyPair = () => {
    return crypto.generateKeyPairSync('rsa', {
        modulusLength: 2048,
        publicKeyEncoding: { type: 'spki', format: 'pem' },
        privateKeyEncoding: { type: 'pkcs8', format: 'pem' },
    });
};

const addUser = async () => {
    const { privateKey: privKey, publicKey: pubKey } = generateRSAKeyPair();
    aid = crypto.randomUUID();
    const response = await fetch(`${serverUrl}/api/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            aid: aid,
            publicKey: pubKey,
            ip: '127.0.2.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    privateKey = privKey;
    publicKey = pubKey;
};

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
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            aid: aid,
            sign: sign.toString('base64'),
            timestamp: timestamp,
            ip: '127.0.2.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
    return data.content; // Return token for Ask test
}

async function testRegisterAlias() {
    console.log('Testing Register Alias...');
    const response = await fetch(`${serverUrl}/usage/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            username: userName,
            password: 'testpass',
            ip: '127.0.2.1',
            fingerprint: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.uuid);
    return data.uuid;
}

async function testLoginAliasNoPreLoginError() {
    console.log('Testing Login Alias - No Pre Login and No online...');
    const response = await fetch(`${serverUrl}/usage/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            username: userName,
            password: 'testpass',
            ip: '127.1.2.1',
            fingerprint: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 400);
    assert.strictEqual(data.result, false);
    assert.strictEqual(data.message, 'no online alias');
}

async function testLoginAliasNoPreLoginSuccess() {
    console.log('Testing Login Alias - No Pre Login but online...');
    const response = await fetch(`${serverUrl}/usage/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            username: userName,
            password: 'testpass',
            ip: '127.1.2.1',
            fingerprint: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
}

async function testLoginAliasWithWrongPreLogin() {
    console.log('Testing Login Alias - Wrong Pre Login...');
    const response = await fetch(`${serverUrl}/usage/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            username: userName,
            password: 'testpass',
            ip: '127.2.2.6',
            fingerprint: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 400);
    assert.strictEqual(data.result, false);
    assert.strictEqual(data.message, 'pre login not match');
}

async function testLoginAliasWithPreLogin() {
    console.log('Testing Login Alias - With Pre Login...');
    const response = await fetch(`${serverUrl}/usage/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            username: userName,
            password: 'testpass',
            ip: '127.1.2.1',
            fingerprint: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.strictEqual(5, data.uuid.split('-').length);
}

async function MFALogin() {
    const timestamp = Date.now().toString();
    const sign = crypto.sign('sha256', Buffer.from(timestamp), {
        key: privateKey,
        padding: crypto.constants.RSA_PKCS1_PADDING,
    });
    const response = await fetch(`${serverUrl}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            aid: aid,
            sign: sign.toString('base64'),
            timestamp: timestamp,
            ip: '127.1.2.1',
            browser: 'Chrome',
        }),
    });
    const data = await response.json();
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
}

async function runTests() {
    try {
        startServer();
        await testLogin();
        await testRegisterAlias();
        await testLoginAliasNoPreLoginError();
        await MFALogin();
        await testLoginAliasNoPreLoginSuccess();
        await testLoginAliasWithWrongPreLogin();
        await testLoginAliasWithPreLogin();
        console.log('All tests passed successfully!');
    } catch (error) {
        console.error('Test failed:', error);
    } finally {
        stopServer();
    }
}

runTests().then(() =>
    console.log('Tests completed')
).catch(e =>
    console.error('Tests failed:', e)
);