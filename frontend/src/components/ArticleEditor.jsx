import { useState, useContext } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import * as Yup from 'yup';
import { createArticle, uploadImage } from '../config/api';
import { AuthContext } from '../context/authContext';
import Navbar from './Navbar';
import { Navigate } from 'react-router-dom';
import ReactQuill from 'react-quill';  // اضافه کردن
import 'react-quill/dist/quill.snow.css';  // استایل‌ها (برای تم snow)

function ArticleEditor() {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { user } = useContext(AuthContext);
  const [title, setTitle] = useState('');
  const [content, setContent] = useState('');  // حالا HTML string هست
  const [imageUrl, setImageUrl] = useState('');
  const [imageFile, setImageFile] = useState(null);

  const schema = Yup.object().shape({
    title: Yup.string().required(t('error.required')).min(3, t('error.minLength', { count: 3 })),
    content: Yup.string().required(t('error.required')).min(10, t('error.minLength', { count: 10 })),
    imageUrl: Yup.string().url(t('error.invalidUrl')).nullable(),
  });

  const modules = {  // تنظیم toolbar برای ویرایشگر
    toolbar: [
      [{ 'header': [1, 2, 3, false] }],
      ['bold', 'italic', 'underline', 'strike', 'blockquote'],
      [{ 'list': 'ordered' }, { 'list': 'bullet' }, { 'indent': '-1' }, { 'indent': '+1' }],
      ['link', 'image'],
      ['clean']
    ],
  };

  const formats = [  // فرمت‌های مجاز
    'header',
    'bold', 'italic', 'underline', 'strike', 'blockquote',
    'list', 'bullet', 'indent',
    'link', 'image'
  ];

  const handleSubmit = async () => {
    try {
      await schema.validate({ title, content, imageUrl }, { abortEarly: false });
      let finalImageUrl = imageUrl;
      if (imageFile) {
        const formData = new FormData();
        formData.append('image', imageFile);
        const response = await uploadImage(formData);
        finalImageUrl = response.imageUrl;
      }
      await createArticle({ title, content, imageUrl: finalImageUrl, type: 'article', lang: i18n.language });
      toast.success(t('submit') + ' موفق');
      navigate('/');
    } catch (error) {
      toast.error(error.response?.data?.error || 'خطا در ارسال');
    }
  };

  if (!user) return <Navigate to="/login" />;

  return (
    <div className={`min-h-screen ${i18n.language === 'fa' ? 'direction-rtl' : 'direction-ltr'}`}>
      <Navbar />
      <div className="container mx-auto p-4">
        <h2 className="text-2xl font-bold mb-4">{t('article')}</h2>
        <input
          className="w-full p-2 border rounded mb-4"
          value={title}
          onChange={e => setTitle(e.target.value)}
          placeholder={t('title')}
        />
        <ReactQuill
          theme="snow"
          value={content}
          onChange={setContent}
          modules={modules}
          formats={formats}
          placeholder={t('content')}
          className="mb-4"
        />
        <input
          className="w-full p-2 border rounded mb-4"
          value={imageUrl}
          onChange={e => setImageUrl(e.target.value)}
          placeholder={t('imageUrl')}
        />
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
          className="bg-blue-500 text-white px-4 py-2 rounded"
          onClick={handleSubmit}
        >
          {t('submit')}
        </button>
      </div>
    </div>
  );
}

export default ArticleEditor;