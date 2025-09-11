import axios from 'axios';
import Cookies from 'js-cookie';

const api = axios.create({
    baseURL: import.meta.env.VITE_API_URL,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json',
    },
});

api.interceptors.request.use(config => {
    const token = Cookies.get('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
});

api.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            Cookies.remove('token');
            window.location.href = '/login';
        }
        // اطمینان از اینکه error.message همیشه رشته است
        const errorMessage = error.response?.data?.error || error.message || 'Unknown error';
        return Promise.reject({ ...error, message: errorMessage });
    }
);

export const loginUser = async (credentials) => {
  const { data } = await api.post('/login', credentials);
  // The server now sets a HTTP-only cookie, so we do not get a token in the response body.
  // The token will be set in the cookie and the frontend just receives a success message.
  return data;
};

export const registerUser = async (credentials) => {
  const { data } = await api.post('/register', credentials);
  // Same as login, the server sets the token in an HTTP-only cookie.
  return data;
};

export const getPosts = async ({ queryKey }) => {
    const [, lang, type, page = 1, limit = 10] = queryKey;
    const { data } = await api.get(`/posts?lang=${lang}&type=${type}&page=${page}&limit=${limit}`);
    return data.data; // { posts: [...] }
};

export const getPostById = async (id) => {
    const { data } = await api.get(`/posts/${id}`);
    return data.data; // PostResponse
};

export const createPost = async (post) => {
  const { data } = await api.post('/posts', post);
  return data;
};

export const createArticle = async (article) => {
  const { data } = await api.post('/articles', article);
  return data;
};

export const uploadImage = async (formData) => {
  const { data } = await api.post('/upload-image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  return data;
};

export const getUserProfile = async (username) => {
  const { data } = await api.get(`/users/${username}`);
  return data;
};


export const getUserPosts = async (username) => {
  const { data } = await api.get(`/users/${username}/posts`);
  return data;
};

export const likePost = async (postId) => {
  const { data } = await api.post(`/posts/${postId}/like`);
  return data;
};

export const unlikePost = async (postId) => {
  const { data } = await api.delete(`/posts/${postId}/unlike`);
  return data;
};
