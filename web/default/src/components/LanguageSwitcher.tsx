import React from 'react';
import { useTranslation } from 'react-i18next';

interface LanguageSwitcherProps {
  className?: string;
}

export function LanguageSwitcher({ className = '' }: LanguageSwitcherProps) {
  const { i18n } = useTranslation();

  return (
    <div className={`lang-switcher ${className}`}>
      <button
        className={`lang-btn ${i18n.language === 'zh' ? 'active' : ''}`}
        onClick={() => i18n.changeLanguage('zh')}
      >
        中文
      </button>
      <button
        className={`lang-btn ${i18n.language === 'en' ? 'active' : ''}`}
        onClick={() => i18n.changeLanguage('en')}
      >
        EN
      </button>
    </div>
  );
}
