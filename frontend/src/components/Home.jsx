import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { getPosts } from '../config/api';
import { useAuth } from '../context/authHooks';
import { useTranslation } from 'react-i18next';
import Navbar from './Navbar';
import PostCard from './ui/PostCard';
import LoadingSpinner from './ui/LoadingSpinner';

const Home = () => {
  const { t, i18n } = useTranslation();
  // `user` در اینجا برای بررسی وضعیت ورود استفاده می‌شود، اما ریدایرکتی صورت نمی‌گیرد.
  const { user } = useAuth(); 

  const [currentLang, setCurrentLang] = React.useState(i18n.language);
  const [postType, setPostType] = React.useState('post');

  const {
    data,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['posts', currentLang, postType],
    queryFn: () => getPosts({ lang: currentLang, type: postType }),
  });

  // اطمینان حاصل می‌کنیم که data.data.posts همیشه یک آرایه است.
  const postsData = data?.data?.posts || [];

  const handleLangChange = (lang) => {
    i18n.changeLanguage(lang);
    setCurrentLang(lang);
  };

  if (isLoading) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <LoadingSpinner />
      </div>
    );
  }

  if (isError) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="container mx-auto mt-10 p-4 text-center">
          <h2 className="text-2xl font-bold text-red-600">{t('error.fetchPosts')}</h2>
          <p className="mt-2 text-gray-600">{error.message}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 dark:bg-gray-900">
      <Navbar />
      <div className="container mx-auto p-4 md:p-8">
        <div className="flex flex-col md:flex-row justify-between items-start md:items-center mb-6">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-4 md:mb-0">
            {t('recent_posts')}
          </h1>
          <div className="flex items-center space-x-4">
            <select
              value={currentLang}
              onChange={(e) => handleLangChange(e.target.value)}
              className="p-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="fa">فارسی</option>
              <option value="en">English</option>
            </select>
            <select
              value={postType}
              onChange={(e) => setPostType(e.target.value)}
              className="p-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
            >
              <option value="post">Post</option>
              <option value="article">Article</option>
            </select>
          </div>
        </div>

        {postsData.length > 0 ? (
          <div className="grid gap-6 sm:grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4">
            {postsData.map((post) => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        ) : (
          <div className="flex flex-col items-center justify-center p-10 text-center">
            <h2 className="text-xl font-semibold text-gray-500 dark:text-gray-400">
              {t('noPostsFound')}
            </h2>
            <p className="mt-2 text-gray-400 dark:text-gray-500">{t('tryCreatingOne')}</p>
          </div>
        )}
      </div>
    </div>
  );
};

export default Home;