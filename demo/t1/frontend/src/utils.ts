import {Aid, AidList, AidType} from "aid-js-sdk"

export const readAidListFromLocalStorage = (): AidList => {
    const defaultUserInfosZip = localStorage.getItem('defaultUserInfosZip');
    const aidsZip = localStorage.getItem('aidsZip');
    if (defaultUserInfosZip === null || aidsZip === null) {
        console.log("No data in local storage");
        return AidList.newAidList();
    }
    return new AidList(defaultUserInfosZip, aidsZip);
}

export const writeAidListToLocalStorage = (aidList: AidList): void => {
    const {
        defaultUserInfosZip,
        aidsZip
    } = aidList.export();

    localStorage.setItem('defaultUserInfosZip', defaultUserInfosZip);
    localStorage.setItem('aidsZip', aidsZip);
}

export const readAid = (aid: string): Aid | null => {
    const aidStr = localStorage.getItem(aid);
    if (aidStr === null) {
        return null;
    }
    return Aid.fromStr(aidStr);
}

export const writeAid = (aid: Aid): void => {
    localStorage.setItem(aid.aid, aid.toStr());
}

export const generateNewAid = async (): Promise<Aid> => {
    const uuid = crypto.randomUUID();
    const newAid = new Aid(uuid, new Map(), new Map(), new Map());
    // rsa generate key pair
    const pair = await window.crypto.subtle.generateKey(
        {
            name: "RSASSA-PKCS1-v1_5",
            modulusLength: 2048, //can be 1024, 2048, or 4096
            publicExponent: new Uint8Array([0x01, 0x00, 0x01]),
            hash: {name: "SHA-256"}, //can be "SHA-1", "SHA-256", "SHA-384", or "SHA-512"
        },
        true,
        ["sign", "verify"]
    )
   // use pem format to store key pair
    const [
        publicKeyPem,
        privateKeyPem
    ] = await Promise.all([
        publicKeyToPem(pair.publicKey),
        privateKeyToPem(pair.privateKey)
    ]);

    newAid.addCert({
        BlockChainUrl: "", ContractAddress: "", ServerAddress: "http://localhost:8080",
        Aid: uuid,
        CertType: AidType.P2p,
        Claims: {},
        Setting: {},
        VerifyOptions: {
            "rsa": publicKeyPem
        }
    }, privateKeyPem);
    writeAid(newAid);
    return newAid;
}

export const getDefaultAid = (aidList: AidList): Aid | undefined => {
    aidList = readAidListFromLocalStorage();
    if (aidList.aids.length === 0) {
        return undefined
    }
    const targetAid = aidList.aids[0];
    let aid = readAid(targetAid.aid);
    if (aid === null) {
        aid = new Aid(targetAid.aid, new Map(), new Map(), new Map());
    }
    return aid;
}

function arrayBufferToBase64(buffer: ArrayBuffer) {
    const bytes = new Uint8Array(buffer);
    let binary = '';
    for (let i = 0; i < bytes.byteLength; i++) {
        binary += String.fromCharCode(bytes[i]);
    }
    return window.btoa(binary);
}

// Function to convert a public key to PEM format
export async function publicKeyToPem(publicKey: CryptoKey): Promise<string> {
    const exported = await window.crypto.subtle.exportKey(
        "spki",
        publicKey
    );
    const exportedAsBase64 = arrayBufferToBase64(exported);
    return `-----BEGIN PUBLIC KEY-----\n${exportedAsBase64}\n-----END PUBLIC KEY-----`;
}


// Function to convert a private key to PEM format
export async function privateKeyToPem(privateKey: CryptoKey): Promise<string> {
    const exported = await window.crypto.subtle.exportKey(
        "pkcs8",
        privateKey
    );
    const exportedAsBase64 = arrayBufferToBase64(exported);
    return `-----BEGIN PRIVATE KEY-----\n${exportedAsBase64}\n-----END PRIVATE KEY-----`;
}

// Function to convert a PEM format private key to CryptoKey
export async function pemToPrivateKey(pemKey: string): Promise<CryptoKey> {
    // Remove the PEM header and footer
    const pemContents = pemKey.replace(
        /(-----BEGIN PRIVATE KEY-----|-----END PRIVATE KEY-----|\s)/g,
        ''
    );

    // Convert base64 to ArrayBuffer
    const binaryDer = base64ToArrayBuffer(pemContents);

    // Import the key
    return await window.crypto.subtle.importKey(
        "pkcs8",
        binaryDer,
        {
            name: "RSASSA-PKCS1-v1_5",
            hash: {name: "SHA-256"}, //can be "SHA-1", "SHA-256", "SHA-384", or "SHA-512"
        },
        true,
        ["sign"]
    );
}

// Helper function to convert base64 to ArrayBuffer
function base64ToArrayBuffer(base64: string): ArrayBuffer {
    const binaryString = window.atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
        bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
}