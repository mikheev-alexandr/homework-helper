import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { signUpTeacher } from '../../services/api.ts';
import '../../styles/SignUp.css';

export default function SignUpTeacher() {
  const { register, handleSubmit } = useForm();
  const [message, setMessage] = useState('');
  const [messageType, setMessageType] = useState('');

  const onSubmit = async (data: any) => {
    setMessage('');
    setMessageType('');
    try {
      await signUpTeacher(data);
      setMessage('Вам на почту отправлено сообщение. Необходимо подтвердить регистрацию.');
      setMessageType('success');
    } catch (error: any) {
      if (error.response) {
        if (error.response.status === 400) {
          setMessage('Пользователь с таким email уже существует');
        } else if (error.response.status === 500) {
          setMessage('Ошибка на сервере. Попробуйте позже.');
        } else {
          setMessage(`Ошибка: ${error.response.data.message || 'Неизвестная ошибка'}`);
        }
      } else {
        setMessage('Сетевая ошибка. Проверьте подключение к интернету.');
      }
      setMessageType('error');
    }
  };

  return (
    <div className="signup-container">
      <form onSubmit={handleSubmit(onSubmit)} className="signup-form">
        <h2>Регистрация для учителей</h2>
        {message && (
          <div className={`message-box ${messageType}`}>
            <p className="message-text">{message}</p>
          </div>
        )}
        <input
          {...register('name')}
          type="name"
          placeholder="ФИО"
          required
          className="form-input"
        />
        <input
          {...register('email')}
          type="email"
          placeholder="Email"
          required
          className="form-input"
        />
        <input
          {...register('password')}
          type="password"
          placeholder="Пароль"
          required
          className="form-input"
        />
        <button type="submit" className="form-button">Зарегистрироваться</button>
      </form>
    </div>
  );
}
