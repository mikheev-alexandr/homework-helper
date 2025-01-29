import React, { useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { resetPassword } from '../../services/api.ts';

export default function UpdatePassword() {
  const [password, setPassword] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [searchParams] = useSearchParams();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setMessage('');
    const token = searchParams.get('token');
    if (!token) {
      setError('Токен отсутствует или недействителен.');
      return;
    }

    try {
      await resetPassword(token, { password });
      setMessage('Пароль успешно обновлён.');
    } catch (err: any) {
      setError('Не удалось обновить пароль. Проверьте токен.');
    }
  };

  return (
    <div className='signin-container'>
      <form onSubmit={handleSubmit} className='signin-form'>
      <h2>Установите новый пароль</h2>
        <input className='form-input'
          type="password"
          placeholder="Введите новый пароль"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
        <button type="submit" className='form-button'>Сохранить</button>
      </form>
      {message && <p>{message}</p>}
      {error && <p style={{ color: 'red' }}>{error}</p>}
    </div>
  );
}
