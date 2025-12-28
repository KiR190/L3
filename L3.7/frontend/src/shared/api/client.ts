import axios from 'axios';

const apiClient = axios.create({
  // Prefer relative baseURL so it works both:
  // - behind nginx at http://localhost:3000 (where /api is proxied to backend)
  // - in Vite dev with proxy (see vite.config.ts)
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor to add JWT token
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('auth_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Response interceptor to handle 401 errors
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    const status = error.response?.status;
    const url: string | undefined = error.config?.url;

    // Do not force-redirect for failed login/register attempts; let the form show an error.
    const isAuthRequest = typeof url === 'string' && (url.includes('/auth/login') || url.includes('/auth/register'));

    if (status === 401 && !isAuthRequest) {
      localStorage.removeItem('auth_token');
      window.location.href = '/auth/login';
    }
    return Promise.reject(error);
  }
);

export default apiClient;

