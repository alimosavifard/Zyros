import { useFormik } from 'formik';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { toast } from 'react-toastify';
import * as Yup from 'yup';
import { registerUser } from '../config/api';
import { useAuth } from '../context/authHooks';
import Navbar from './Navbar';

const Register = () => {
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const { login } = useAuth();

  const schema = Yup.object().shape({
    username: Yup.string().required(t('error.required')).min(3, t('error.minLength', { count: 3 })),
    password: Yup.string().required(t('error.required')).min(6, t('error.minLength', { count: 6 })),
  });

  const formik = useFormik({
    initialValues: { username: '', password: '' },
    validationSchema: schema,
    onSubmit: async (values, { setSubmitting, setErrors }) => {
      try {
        await registerUser(values);  // بدون دریافت token
        await login();  // فراخوانی بدون آرگومان، از کوکی می‌خونه
        toast.success(t('register') + ' موفق');
        navigate('/');
      } catch (error) {
        if (error.response?.status === 409) {  // فرض بر این که بک‌اند 409 برای duplicate برمی‌گردونه
          setErrors({ username: t('error.duplicateUsername') });
        } else {
          toast.error('خطا در ثبت‌نام');
        }
      } finally {
        setSubmitting(false);
      }
    },
  });

  return (
    <div className={`min-h-screen ${i18n.language === 'fa' ? 'direction-rtl' : 'direction-ltr'}`}>
      <Navbar />
      <form onSubmit={formik.handleSubmit} className="max-w-md mx-auto mt-10 space-y-4">
        <div>
          <input
            name="username"
            placeholder={t('username')}
            value={formik.values.username}
            onChange={formik.handleChange}
            className="w-full border p-2 rounded"
          />
          {formik.errors.username && <div className="text-red-500">{formik.errors.username}</div>}
        </div>
        <div>
          <input
            name="password"
            type="password"
            placeholder={t('password')}
            value={formik.values.password}
            onChange={formik.handleChange}
            className="w-full border p-2 rounded"
          />
          {formik.errors.password && <div className="text-red-500">{formik.errors.password}</div>}
        </div>
        <button
          type="submit"
          disabled={formik.isSubmitting}
          className="bg-blue-600 text-white px-4 py-2 rounded disabled:opacity-50"
        >
          {t('register')}
        </button>
      </form>
    </div>
  );
};

export default Register;