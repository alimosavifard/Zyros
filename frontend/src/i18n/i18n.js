import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';
import HttpBackend from 'i18next-http-backend';

i18n
  .use(HttpBackend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: 'fa',
    supportedLngs: ['fa', 'en'],
    backend: {
      loadPath: '/assets/locales/{{lng}}/translation.json', // مسیر اصلاح شده
      addPath: '/assets/locales/{{lng}}/missing.json',
      cache: {
        enabled: true,
        expirationTime: 7 * 24 * 60 * 60 * 1000,
      },
    },
    interpolation: { escapeValue: false },
    detection: {
      order: ['cookie', 'navigator', 'htmlTag'],
      caches: ['cookie'],
      cookieMinutes: 60 * 24 * 7,
    },
  });

export default i18n;