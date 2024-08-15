interface AIAuthRequest {
  username?: string;
  password?: string;
  token?: string;
  ip?: string;
  fingerPrint?: string;
}

interface AIAuthResponse {
  token?: string;
  uuid?: string;
  message?: string;
  result?: boolean;
}

// 通用的 fetch 錯誤處理函數
async function handleResponse(response: Response) {
  if (!response.ok) {
    const errorBody = await response.text();
    throw new Error(
      `HTTP error! status: ${response.status}, body: ${errorBody}`,
    );
  }
  return response.json();
}

// API 請求函數
export const apiService = {
  login: async (request: AIAuthRequest): Promise<AIAuthResponse> => {
    try {
      const response = await fetch(`/api/auth/login`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });
      return handleResponse(response);
    } catch (error) {
      console.error("Login error:", error);
      throw error;
    }
  },

  register: async (request: AIAuthRequest): Promise<AIAuthResponse> => {
    try {
      const response = await fetch(`/api/auth/register`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(request),
      });
      return handleResponse(response);
    } catch (error) {
      console.error("Register error:", error);
      throw error;
    }
  },
};

// 使用示例：
// import { apiService } from './api-service';
//
// const loginUser = async () => {
//   try {
//     const response = await apiService.login({
//       username: 'user',
//       password: 'pass',
//       loginType: 'password',
//       rememberMe: true
//     });
//     console.log(response.uuid);
//   } catch (error) {
//     console.error('Login failed:', error);
//   }
// };
