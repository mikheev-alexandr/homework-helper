import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { signInTeacher, signInStudent } from '../../services/api.ts';
import '../../styles/SignIn.css';

interface SignInProps {
  role: 'teacher' | 'student';
}

export default function SignIn({ role }: SignInProps) {
  const { register, handleSubmit } = useForm();
  const [errorMessage, setErrorMessage] = useState('');
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  const onSubmit = async (data: any) => {
    setErrorMessage('');
    try {
      const response =
        role === 'teacher'
          ? await signInTeacher({ email: data.email, password: data.password })
          : await signInStudent({ code: data.code, password: data.password });

      const redirect =
        searchParams.get('redirect') ||
        (role === 'teacher' ? '/teacher/students' : '/student/teachers');

      navigate(redirect);
    } catch (error: any) {
      if (error.response) {
        if (error.response.status === 401) {
          setErrorMessage(
            role === 'teacher'
              ? 'Неправильный email или пароль.'
              : 'Неправильное кодовое слово или пароль.'
          );
        } else if (error.response.status === 500) {
          setErrorMessage('Ошибка на сервере. Попробуйте позже.');
        } else {
          setErrorMessage(`Ошибка: ${error.response.data.message || 'Неизвестная ошибка'}`);
        }
      } else {
        setErrorMessage('Сетевая ошибка. Проверьте подключение к интернету.');
      }
    }
  };

  return (
    <div className="signin-container">
      <form onSubmit={handleSubmit(onSubmit)} className="signin-form">
        <h2>Вход для {role === 'teacher' ? 'учителей' : 'учеников'}</h2>
        {errorMessage && (
          <div className="error-box">
            <p className="error-message">{errorMessage}</p>
          </div>
        )}
        {role === 'teacher' ? (
          <>
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
            <div className="forgot-password-link">
              <Link to="/auth/teacher/reset-password">
                Забыли пароль?
              </Link>
            </div>
          </>
        ) : (
          <>
            <input
              {...register('code')}
              type="text"
              placeholder="Кодовое слово"
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
          </>
        )}
        <button type="submit" className="form-button">
          Войти
        </button>
      </form>
    </div>
  );
}
