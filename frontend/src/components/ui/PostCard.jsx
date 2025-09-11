import React from 'react';
import { Link } from 'react-router-dom';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { toast } from 'react-toastify';
import { likePost, unlikePost } from '../../config/api';
import { useAuth } from '../../context/authHooks';
import { formatDistanceToNow } from 'date-fns';
import { faIR } from 'date-fns/locale';

const PostCard = ({ post }) => {
  const { t, i18n } = useTranslation();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  
  // A hypothetical `isLiked` property could be included in the API response
	const isLiked = post.isLikedByUser ?? false; // fallback به false
	const likesCount = post.likesCount ?? 0; // fallback به 0


  const likeMutation = useMutation({
    mutationFn: isLiked ? unlikePost : likePost,
    onSuccess: () => {
      // Invalidate the posts query to refetch and update like counts
      queryClient.invalidateQueries(['posts']);
      queryClient.invalidateQueries(['post', post.id]); // Also invalidate single post view
      toast.success(isLiked ? t('post.unliked') : t('post.liked'));
    },
    onError: (error) => {
      toast.error(t('error.failedToLike'));
      console.error("Liking failed:", error);
    },
  });

  const handleLike = () => {
    if (!user) {
      toast.info(t('auth.loginRequired'));
      return;
    }
    likeMutation.mutate(post.id);
  };
  
  const imageUrl = post.imageUrl || `https://source.unsplash.com/random/400x200?sig=${post.id}`;

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md hover:shadow-xl transition-shadow duration-300 overflow-hidden h-full flex flex-col">
      <Link to={`/posts/${post.id}`} className="block">
        <img 
          src={imageUrl} 
          alt={post.title} 
          className="w-full h-48 object-cover transition-transform duration-300 hover:scale-105"
        />
      </Link>
      <div className="p-4 flex flex-col flex-grow">
        <h3 className="text-xl font-semibold text-gray-900 dark:text-white leading-tight mb-2">
          {post.title}
        </h3>
        <p className="text-gray-600 dark:text-gray-400 text-sm mb-3 flex-grow line-clamp-3">
          {post.content}
        </p>
        
        <div className="flex justify-between items-center text-xs text-gray-500 dark:text-gray-400 mt-auto pt-2 border-t border-gray-100 dark:border-gray-700">
          <Link to={`/users/${post.user.username}`} className="hover:underline">
            <span>{t('by')} {post.user.username}</span>
          </Link>
          <span>
            {formatDistanceToNow(new Date(post.createdAt), {
              addSuffix: true,
              locale: i18n.language === 'fa' ? fa : undefined,
            })}
          </span>
        </div>
        
        <div className="flex items-center space-x-2 mt-2">
          <button 
            onClick={handleLike} 
            className={`flex items-center space-x-1 ${isLiked ? 'text-red-500' : 'text-gray-500 dark:text-gray-400'} hover:text-red-600 transition-colors`}
            disabled={likeMutation.isLoading}
          >
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              className="h-5 w-5 fill-current" 
              viewBox="0 0 20 20" 
              fill="currentColor"
            >
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.381-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
            <span>{likesCount}</span>
          </button>
        </div>
      </div>
    </div>
  );
};

export default PostCard;