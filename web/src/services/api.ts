import axios, { AxiosError } from 'axios';

const apiClient = axios.create({
  baseURL: '/api/v1',
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Paste {
  id: string;
  content: string;
  language: string;
  created_at: string;
  is_protected?: boolean; // Новое поле
  burn_after_reading?: boolean; // Новое поле
}

// Интерфейс для настроек создания
export interface CreateOptions {
  content: string;
  language: string;
  ttl: number;      // Секунды
  password?: string;
  burn?: boolean;
}

export default {
  async createPaste(options: CreateOptions) {
    const response = await apiClient.post<Paste>('/paste', options);
    return response.data;
  },

  async getPaste(id: string, password?: string) {
    try {
      const config = password ? { params: { password } } : {};
      const response = await apiClient.get<Paste>(`/paste/${id}`, config);
      return response.data;
    } catch (error) {
      // Если сервер вернул 403 (Password Required), прокидываем ошибку со специальным полем
      if (axios.isAxiosError(error) && error.response?.status === 403) {
        throw { status: 403, isProtected: true };
      }
      throw error;
    }
  },
};