import axios from 'axios';

// Создаем инстанс с базовыми настройками
const apiClient = axios.create({
  baseURL: 'http://localhost:8080/api/v1', // Адрес нашего Go сервера
  headers: {
    'Content-Type': 'application/json',
  },
});

export interface Paste {
  id: string;
  content: string;
  language: string;
  created_at: string;
}

export default {
  // Создать пасту
  async createPaste(content: string, language: string = 'plaintext') {
    const response = await apiClient.post<Paste>('/paste', {
      content,
      language,
    });
    return response.data;
  },

  // Получить пасту
  async getPaste(id: string) {
    const response = await apiClient.get<Paste>(`/paste/${id}`);
    return response.data;
  },
};