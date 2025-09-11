// src/components/UserProfile.jsx

import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { useParams, Navigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { getUserProfile, getUserPosts } from '../config/api';
import { useAuth } from '../context/authHooks';
import Navbar from './Navbar';
import LoadingSpinner from './ui/LoadingSpinner';
import PostCard from './ui/PostCard';

const UserProfile = () => {
  const { username } = useParams();
  const { t } = useTranslation();
  const { user } = useAuth();

  const { data: profileData, isLoading: profileLoading, isError: profileError } = useQuery({
    queryKey: ['userProfile', username],
    queryFn: () => getUserProfile(username),
    enabled: !!username,
  });

  const { data: postsData, isLoading: postsLoading, isError: postsError } = useQuery({
    queryKey: ['userPosts', username],
    queryFn: () => getUserPosts(username),
    enabled: !!username,
  });

  if (!user) {
    return <Navigate to="/login" />;
  }

  if (profileLoading || postsLoading) {
    return <LoadingSpinner />;
  }

  if (profileError || postsError) {
    return (
      <div className="text-center mt-10 text-red-500">
        {t('error.failedToFetchUser')}
      </div>
    );
  }

  const profile = profileData.data;
  const posts = postsData.data.posts;

  return (
    <div className="bg-gray-100 min-h-screen">
      <Navbar />
      <div className="container mx-auto p-4 max-w-4xl">
        {/* User Profile Card */}
        <div className="bg-white rounded-lg shadow-md p-6 mb-8 text-center">
          <img
            src={profile.profilePictureUrl || 'https://via.placeholder.com/150'}
            alt={`${profile.username}'s profile`}
            className="w-32 h-32 rounded-full mx-auto mb-4 border-4 border-gray-300"
          />
          <h1 className="text-3xl font-bold text-gray-900">{profile.username}</h1>
          {profile.bio && <p className="text-gray-600 mt-2">{profile.bio}</p>}
        </div>

        {/* User Posts Section */}
        <h2 className="text-2xl font-bold text-gray-800 mb-4">{t('userPosts')}</h2>
        {posts && posts.length > 0 ? (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {posts.map(post => (
              <PostCard key={post.id} post={post} />
            ))}
          </div>
        ) : (
          <div className="text-center text-gray-500 mt-10">{t('noPostsFound')}</div>
        )}
      </div>
    </div>
  );
};

export default UserProfile;