// components/Navbar.jsx
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/authHooks';
import { toast } from 'react-toastify';

function Navbar() {
  const { t, i18n } = useTranslation();
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const changeLanguage = (lng) => {
    i18n.changeLanguage(lng);
  };
  
  const handleLogout = () => {
    logout();
    toast.success(t('logout_success'));
    navigate('/login');
  };

  return (
    <nav className={`bg-gray-800 text-white p-4 ${i18n.language === 'fa' ? 'direction-rtl' : 'direction-ltr'}`}>
      <div className="container mx-auto flex justify-between items-center">
        
        {/* Section 1: Main Navigation (Always visible) */}
        <div className="flex items-center space-x-4">
          <Link to="/" className="hover:underline text-lg font-bold">
            {t('welcome')}
          </Link>
        </div>

        {/* Section 2: Conditional User Links */}
        <div className="flex items-center space-x-4">
          {user ? (
            // If user is logged in
            <>
              <Link to={`/users/${user.username}`} className="hover:underline">
                {t('profile')}
              </Link>
              
              <Link to="/post" className="hover:underline">
                {t('create_post')}
              </Link>
              <Link to="/article" className="hover:underline">
                {t('create_article')}
              </Link>

              <button onClick={handleLogout} className="hover:underline">
                {t('logout')}
              </button>
            </>
          ) : (
            // If user is not logged in
            <>
              <Link to="/login" className="hover:underline">
                {t('login')}
              </Link>
              <Link to="/register" className="hover:underline">
                {t('register')}
              </Link>
            </>
          )}

          {/* Section 3: Language Switcher */}
          <div className="flex items-center space-x-2">
            <button onClick={() => changeLanguage('fa')} className="hover:underline">
              فارسی
            </button>
            <button onClick={() => changeLanguage('en')} className="hover:underline">
              English
            </button>
          </div>
        </div>
      </div>
    </nav>
  );
}

export default Navbar;