import flourite from 'flourite';
import { supportedLanguages } from '../config/languages';

export const detectLanguage = (content: string): string | null => {
  // Не пытаемся определить, если контент слишком короткий
  if (!content || content.trim().length < 10) return null;

  // flourite возвращает объект { language: "Javascript", ... }
  const result = flourite(content, { noUnknown: true });
  
  const detected = result.language.toLowerCase();
  
  if (detected === 'unknown') return null;

  // Ищем совпадение в нашем конфиге
  const match = supportedLanguages.find(lang => {
    // Сравниваем с основным именем
    if (lang.name.toLowerCase() === detected) return true;
    // Сравниваем с алиасами
    if (lang.alias.some(a => a.toLowerCase() === detected)) return true;
    return false;
  });

  return match ? match.name : null;
};