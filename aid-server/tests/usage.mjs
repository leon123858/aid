import {execSync, spawn} from 'child_process';
import assert from 'assert';
import crypto from 'crypto';

const serverUrl = 'http://127.0.0.1:8080';
// const serverUrl = 'http://20.2.209.109';
let serverProcess;
let userName, publicKey, privateKey, aid;

const startServer = () => {
    serverProcess = spawn('../bin/aid', ['server']);
    execSync('sleep 2');
};

const stopServer = () => {
    if (serverProcess) {
        serverProcess.kill();
    }
};

const generateRSAKeyPair = () => {
    return crypto.generateKeyPairSync('rsa', {
        modulusLength: 2048,
        publicKeyEncoding: {type: 'spki', format: 'pem'},
        privateKeyEncoding: {type: 'pkcs8', format: 'pem'},
    });
};

const makeRequest = async (endpoint, method, body) => {
    const response = await fetch(`${serverUrl}${endpoint}`, {
        method,
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify(body),
    });
    return {response, data: await response.json()};
};

const addUser = async () => {
    const {privateKey: privKey, publicKey: pubKey} = generateRSAKeyPair();
    aid = crypto.randomUUID();
    const {response, data} = await makeRequest('/api/register', 'POST', {
        aid,
        publicKey: pubKey,
        ip: '127.0.2.1',
        browser: 'Chrome',
    });
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    privateKey = privKey;
    publicKey = pubKey;
};

const testLogin = async (fingerprint = 'Chrome') => {
    console.log('Testing Login...');
    await addUser();
    const timestamp = Date.now().toString();
    const sign = crypto.sign('sha256', Buffer.from(timestamp), {
        key: privateKey,
        padding: crypto.constants.RSA_PKCS1_PADDING,
    });
    const {response, data} = await makeRequest('/api/login', 'POST', {
        aid,
        sign: sign.toString('base64'),
        timestamp,
        ip: '127.0.2.1',
        browser: fingerprint,
    });
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
    return data.content;
};

const testRegisterAlias = async (fingerprint = 'Chrome') => {
    console.log('Testing Register Alias...');
    const {response, data} = await makeRequest('/usage/register', 'POST', {
        username: userName,
        password: 'testpass',
        ip: '127.0.2.1',
        fingerprint,
    });
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.uuid);
    return data.uuid;
};

const testLoginAlias = async (scenario, ip, fingerprint = 'Chrome') => {
    console.log(`Testing Login Alias - ${scenario}...`);
    const {response, data} = await makeRequest('/usage/login', 'POST', {
        username: userName,
        password: 'testpass',
        ip,
        fingerprint,
    });
    return {response, data};
};

const MFALogin = async (fingerprint = 'Chrome') => {
    const timestamp = Date.now().toString();
    const sign = crypto.sign('sha256', Buffer.from(timestamp), {
        key: privateKey,
        padding: crypto.constants.RSA_PKCS1_PADDING,
    });
    const {response, data} = await makeRequest('/api/login', 'POST', {
        aid,
        sign: sign.toString('base64'),
        timestamp,
        ip: '127.1.2.1',
        browser: fingerprint,
    });
    assert.strictEqual(response.status, 200);
    assert.strictEqual(data.result, true);
    assert.ok(data.content);
};

const runTestSuite = async (browser) => {
    console.log(`Running test suite for ${browser}...`);

    await testLogin(browser);
    await testRegisterAlias(browser);

    let result = await testLoginAlias('No Pre Login and No online', '127.1.2.1', browser);
    assert.strictEqual(result.response.status, 400);
    assert.strictEqual(result.data.result, false);
    assert.strictEqual(result.data.message, 'no online alias');

    await MFALogin(browser);

    result = await testLoginAlias('No Pre Login but online', '127.1.2.1', browser);
    assert.strictEqual(result.response.status, 200);
    assert.strictEqual(result.data.result, true);

    result = await testLoginAlias('Wrong Pre Login', '127.2.2.6', browser);
    assert.strictEqual(result.response.status, 400);
    assert.strictEqual(result.data.result, false);

    result = await testLoginAlias('With Pre Login', '127.1.2.1', browser);
    assert.strictEqual(result.response.status, 200);
    assert.strictEqual(result.data.result, true);
    assert.strictEqual(result.data.uuid.split('-').length, 5);
};

const runAllTests = async () => {
    try {
        startServer();
        userName = crypto.randomUUID();

        console.log('First AID test...');
        await runTestSuite('Chrome');

        console.log('Second AID test...');
        await runTestSuite('firefox');

        console.log('All tests passed successfully!');
    } catch (error) {
        console.error('Test failed:', error);
    } finally {
        stopServer();
    }
};

runAllTests().then(() =>
    console.log('Tests completed')
).catch(e =>
    console.error('Tests failed:', e)
);