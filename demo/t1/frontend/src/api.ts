import { AidCert } from 'aid-js-sdk';

const API_BASE_URL = 'http://localhost:8080';

export interface TodoItem {
    id: number;
    task: string;
    done: boolean;
}

export class TodoApiClient {
    private readonly aid: string | null = null;
    private readonly privateKey: CryptoKey | null = null;
    private readonly cert : AidCert | null = null;

    constructor(aid: string, privateKey: CryptoKey, cert: AidCert) {
        this.aid = aid;
        this.privateKey = privateKey
        this.cert = cert
    }

    async login(cert: AidCert): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');

        const response = await fetch(`${API_BASE_URL}/login/${this.aid}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ cert }),
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        return response.json();
    }

    async logout(): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');

        const { sign, preSign } = await this.generateSignature();
        const response = await fetch(`${API_BASE_URL}/logout/${this.aid}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Sign': sign,
                'PreSign': preSign,
            },
        });

        if (!response.ok) {
            throw new Error('Logout failed');
        }

        return response.json();
    }

    async getTodos(aidStr: string): Promise<TodoItem[]> {
        if (!aidStr) {
            throw new Error('AID is not set');
        }

        const response = await fetch(`${API_BASE_URL}/todos/${aidStr}`, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error('Failed to get todos');
        }

        return response.json();
    }

    async createTodos(todos: TodoItem[]): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');

        const { sign, preSign } = await this.generateSignature();
        const response = await fetch(`${API_BASE_URL}/todos/${this.aid}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Sign': sign,
                'PreSign': preSign,
            },
            body: JSON.stringify(todos),
        });

        if (!response.ok) {
            throw new Error('Failed to create todos');
        }

        return response.json();
    }

    private async generateSignature(): Promise<{ sign: string, preSign: string }> {
        if (!this.privateKey) throw new Error('Private key is not set');

        const preSign = Date.now().toString();
        // sign the hashed preSign
        const signature = await window.crypto.subtle.sign(
            {
                name: "RSASSA-PKCS1-v1_5",
                hash: {name: "SHA-256"}, //can be "SHA-1", "SHA-256", "SHA-384", or "SHA-512"
            },
            this.privateKey,
            new TextEncoder().encode(preSign)
        );
        const sign = btoa(this.uint8ArrayToString(new Uint8Array(signature)));
        return { sign, preSign };
    }

    private uint8ArrayToString(array: Uint8Array): string {
        const chunk = 8192; // 處理大型數組
        let result = '';
        for (let i = 0; i < array.length; i += chunk) {
            result += String.fromCharCode.apply(null, Array.from(array.subarray(i, i + chunk)));
        }
        return result;
    }
}
