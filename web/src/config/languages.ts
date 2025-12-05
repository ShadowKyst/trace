import { javascript } from '@codemirror/lang-javascript';
import { html } from '@codemirror/lang-html';
import { css } from '@codemirror/lang-css';
import { json } from '@codemirror/lang-json';
import { python } from '@codemirror/lang-python';
import { go } from '@codemirror/lang-go';
import { markdown } from '@codemirror/lang-markdown';
import type { LanguageSupport } from '@codemirror/language';

// Определяем тип для нашего списка
interface LanguageOption {
  name: string;
  extension: () => LanguageSupport;
  alias: string[];
}

export const supportedLanguages: LanguageOption[] = [
  { name: 'Plain Text', extension: () => [] as unknown as LanguageSupport, alias: ['txt', 'text'] },
  { name: 'JavaScript', extension: javascript, alias: ['js', 'node'] },
  { name: 'TypeScript', extension: () => javascript({ typescript: true }), alias: ['ts'] },
  { name: 'Go', extension: go, alias: ['golang'] },
  { name: 'Python', extension: python, alias: ['py'] },
  { name: 'JSON', extension: json, alias: ['json'] },
  { name: 'HTML', extension: html, alias: ['html'] },
  { name: 'CSS', extension: css, alias: ['css'] },
  { name: 'Markdown', extension: markdown, alias: ['md'] },
];

export const getExtensionByName = (name: string) => {
  const lang = supportedLanguages.find(l => l.name === name || l.alias.includes(name));
  return lang ? lang.extension() : [];
};