import { useFormik } from 'formik';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { createPost, uploadImage } from '../config/api';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../context/authHooks';
import { toast } from 'react-toastify';
import * as Yup from 'yup';
import Navbar from './Navbar';
import { Navigate } from 'react-router-dom';
import { useState } from 'react';

const PostForm = () => {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const { user } = useAuth();
  const [imageFile, setImageFile] = useState(null);

  const schema = Yup.object().shape({
    title: Yup.string().required(t('error.required')).min(3, t('error.minLength', { count: 3 })),
    content: Yup.string().required(t('error.required')).min(10, t('error.minLength', { count: 10 })),
    imageUrl: Yup.string().url(t('error.invalidUrl')).nullable(), // اختیاری
  });

  const mutation = useMutation({
    mutationFn: async (values) => {
      let imageUrl = values.imageUrl;
      if (imageFile) {
        const formData = new FormData();
        formData.append('image', imageFile);
        const response = await uploadImage(formData);
        imageUrl = response.imageUrl;
      }
      return createPost({ ...values, imageUrl, type: 'post', lang: i18n.language });
    },
    onSuccess: () => {
      queryClient.invalidateQueries(['posts', i18n.language]);
      toast.success(t('submit') + ' موفق');
      navigate('/');
    },
    onError: (error) => {
      toast.error(error.response?.data?.error || 'خطا در ارسال');
    },
  });

  const formik = useFormik({
    initialValues: { title: '', content: '', imageUrl: '' },
    validationSchema: schema,
    onSubmit: (values) => mutation.mutate(values),
  });

  if (!user) return <Navigate to="/login" replace />;

  return (
    <div className={`min-h-screen ${i18n.language === 'fa' ? 'direction-rtl' : 'direction-ltr'}`}>
      <Navbar />
      <form onSubmit={formik.handleSubmit} className="max-w-md mx-auto mt-10 space-y-4">
        <div>
          <input
            name="title"
            placeholder={t('title')}
            value={formik.values.title}
            onChange={formik.handleChange}
            className="w-full border p-2 rounded"
          />
          {formik.errors.title && <div className="text-red-500">{formik.errors.title}</div>}
        </div>
        <div>
          <textarea
            name="content"
            placeholder={t('content')}
            value={formik.values.content}
            onChange={formik.handleChange}
            className="w-full border p-2 rounded h-40"
          />
          {formik.errors.content && <div className="text-red-500">{formik.errors.content}</div>}
        </div>
        <div>
          <input
            name="imageUrl"
            placeholder={t('imageUrl')}
            value={formik.values.imageUrl}
            onChange={formik.handleChange}
            className="w-full border p-2 rounded"
          />
          {formik.errors.imageUrl && <div className="text-red-500">{formik.errors.imageUrl}</div>}
        </div>
        <div>
          <label className="block text-sm font-medium text-gray-700">{t('uploadImage')}</label>
          <input
            type="file"
            accept="image/jpeg,image/png,image/gif"
            onChange={(e) => setImageFile(e.target.files[0])}
            className="w-full border p-2 rounded"
          />
        </div>
        {imageFile && (
          <img
            src={URL.createObjectURL(imageFile)}
            alt="Preview"
            className="mt-4 w-full h-48 object-cover rounded"
          />
        )}
        <button
          type="submit"
          disabled={mutation.isLoading}
          className="bg-green-600 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          {t('submit')}
        </button>
      </form>
    </div>
  );
};

export default PostForm;