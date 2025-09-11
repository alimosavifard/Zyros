import React from 'react';
import { useParams } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { getPostById } from '../config/api';
import Navbar from './Navbar';
import LoadingSpinner from './ui/LoadingSpinner';

const PostDetails = () => {
  const { t } = useTranslation();
  const { id } = useParams(); // دریافت شناسه پست از URL

  // واکشی اطلاعات پست بر اساس شناسه
  const {
    data: post,
    isLoading,
    isError,
    error,
  } = useQuery({
    queryKey: ['post', id],
    queryFn: () => getPostById(id),
    enabled: !!id, // واکشی فقط در صورت وجود شناسه
  });

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
          <h2 className="text-2xl font-bold text-red-600">{t('error.fetchPost')}</h2>
          <p className="mt-2 text-gray-600">{error.message}</p>
        </div>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="min-h-screen">
        <Navbar />
        <div className="container mx-auto mt-10 p-4 text-center">
          <h2 className="text-2xl font-bold text-gray-600">{t('postNotFound')}</h2>
        </div>
      </div>
    );
  }

  // نمایش اطلاعات پست
  return (
    <div className="min-h-screen">
      <Navbar />
      <div className="container mx-auto mt-10 p-4">
        {post.imageUrl && (
          <img
            src={post.imageUrl}
            alt={post.title}
            className="w-full h-96 object-cover rounded-lg shadow-lg mb-6"
          />
        )}
        <h1 className="text-4xl font-bold mb-4">{post.title}</h1>
        <p className="text-gray-600 mb-6">
          {t('lang')}: {post.lang === 'fa' ? t('farsi') : t('english')} | {t('type')}: {t(post.type)}
        </p>
        <div className="prose max-w-none text-gray-800 leading-relaxed dark:text-gray-200">
          <p>{post.content}</p>
        </div>
      </div>
    </div>
  );
};

export default PostDetails;