import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Mengarah ke backend Golang
});

// Otomatis menyisipkan Token di setiap request
api.interceptors.request.use((config) => {
  // Kita pakai localStorage bawaan browser untuk kemudahan
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;