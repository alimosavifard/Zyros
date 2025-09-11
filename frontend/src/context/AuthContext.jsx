import { useEffect, useState, useCallback } from 'react';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';
import PropTypes from 'prop-types';
import { AuthContext } from './authContext';

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isLoading, setIsLoading] = useState(true);

  const decodeToken = useCallback((token) => {
    try {
      const decoded = jwtDecode(token);
      return decoded;
    } catch {
      console.error('Invalid token');
      return null;
    }
  }, []);

  const login = useCallback(async () => {
    const token = Cookies.get('token');
    if (!token) {
      setUser(null);
      setIsLoading(false);
      return;
    }
    const decoded = decodeToken(token);
    if (decoded) {
      setUser({ id: decoded.userID });
    }
    setIsLoading(false);
  }, [decodeToken]);

  const logout = useCallback(() => {
    Cookies.remove('token');
    // The browser will handle the CSRF cookie removal, which is HTTP-only
    setUser(null);
    setIsLoading(false);
  }, []);

  useEffect(() => {
    login();
  }, [login]);

  if (isLoading) {
    return <div>در حال بررسی وضعیت احراز هویت...</div>;
  }

  const value = { user, isLoading, login, logout };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

AuthProvider.propTypes = {
  children: PropTypes.node.isRequired,
};