import React, { useState } from 'react';
import { requestPasswordReset } from '../../services/api.ts';
import '../../styles/SignIn.css'

export default function ResetPassword() {
  const [email, setEmail] = useState('');
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setMessage('');
    try {
      await requestPasswordReset({ email });
      setMessage('Ссылка для сброса пароля отправлена на ваш email.');
    } catch (err: any) {
      setError('Не удалось отправить ссылку. Проверьте email.');
    }
  };

  return (
    <div className='signin-container'>
      <form onSubmit={handleSubmit} className='signin-form'>
        <h2>Сброс пароля</h2>
        {error && (
          <div className="error-box">
            <p className="error-message">{error}</p>
          </div>
        )}
        <input className='form-input'
          type="email"
          placeholder="Введите ваш email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
        <button type="submit" className='form-button'>Отправить</button>
        {message && <p>{message}</p>}
      </form>
    </div>
  );
}
