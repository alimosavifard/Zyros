import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { I18nextProvider } from 'react-i18next';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { ToastContainer } from 'react-toastify';
import { lazy, Suspense } from 'react';
import { AuthProvider } from './context/AuthContext.jsx';
import { useAuth } from './context/authHooks.js';
import i18n from './i18n/i18n';
import 'react-toastify/dist/ReactToastify.css';
import './index.css';

const Home = lazy(() => import('./components/Home'));
const PostForm = lazy(() => import('./components/PostForm'));
const ArticleEditor = lazy(() => import('./components/ArticleEditor'));
const Login = lazy(() => import('./components/Login'));
const Register = lazy(() => import('./components/Register'));
const PostDetails = lazy(() => import('./components/PostDetails'));
const UserProfile = lazy(() => import('./components/UserProfile'));
const Navbar = lazy(() => import('./components/Navbar'));

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      staleTime: 5 * 60 * 1000,
      cacheTime: 10 * 60 * 1000,
    },
  },
});

// New ProtectedRoute component
const ProtectedRoute = ({ children }) => {
  const { user, isLoading } = useAuth();
  if (isLoading) {
    return <div>در حال بررسی...</div>;
  }
  return user ? children : <Navigate to="/login" replace />;
};

function App() {
  return (
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={queryClient}>
        <AuthProvider>
          <Router>
            <Suspense fallback={<div className="text-center mt-10">در حال بارگذاری...</div>}>
              <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/login" element={<Login />} />
                <Route path="/register" element={<Register />} />
                <Route path="/posts/:id" element={<PostDetails />} />

                {/* Protected routes */}
                <Route
                  path="/post"
                  element={<ProtectedRoute><PostForm /></ProtectedRoute>}
                />
                <Route
                  path="/article"
                  element={<ProtectedRoute><ArticleEditor /></ProtectedRoute>}
                />
                <Route
                  path="/users/:username"
                  element={<ProtectedRoute><UserProfile /></ProtectedRoute>}
                />
              </Routes>
            </Suspense>
          </Router>
          <ToastContainer position={i18n.language === 'fa' ? 'top-right' : 'top-left'} />
        </AuthProvider>
      </QueryClientProvider>
    </I18nextProvider>
  );
}

export default App;