import {AidCert, generateSignature} from 'aid-js-sdk';

const API_BASE_URL = 'http://localhost:8080';
const AID_SERVER_URL = 'http://localhost:7001';

export interface TodoItem {
    id: number;
    task: string;
    done: boolean;
}

export class TodoApiClient {
    private readonly aid: string | null = null;
    private readonly privateKey: CryptoKey | null = null;

    constructor(aid: string, privateKey: CryptoKey) {
        this.aid = aid;
        this.privateKey = privateKey
    }

    async login(cert: AidCert): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');

        const response = await fetch(`${API_BASE_URL}/login/${this.aid}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({cert}),
        });

        if (!response.ok) {
            throw new Error('Login failed');
        }

        return response.json();
    }

    async logout(): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');
        if (!this.privateKey) throw new Error('Private key is not set');

        const preSign = new Date().getTime().toString();
        const sign = await generateSignature(this.privateKey, preSign);
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
        if (!this.privateKey) throw new Error('AID is not set');
        const preSign = new Date().getTime().toString();
        const sign = await generateSignature(this.privateKey, preSign);
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

    async registerAidRemote(certHash: string): Promise<{ result: string }> {
        if (!this.aid) throw new Error('AID is not set');
        const response = await fetch(`${AID_SERVER_URL}/register/cert`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({aid: this.aid, hash: certHash}),
        });

        if (!response.ok) {
            throw new Error('Failed to register remote aid');
        }

        return response.json();
    }
}
