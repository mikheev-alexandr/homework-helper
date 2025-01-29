import React from 'react';
import { Link } from 'react-router-dom';
import '../styles/NotFoundPage.css';

export default function NotFoundPage() {
  return (
    <div className="not-found-container">
      <h1>404</h1>
      <p>Упс! Страница не найдена.</p>
      <p>
        Возможно, вы ошиблись в адресе или страница была удалена. Вы можете вернуться на{' '}
        <Link to="/" className="link">главную страницу</Link>.
      </p>
    </div>
  );
}
