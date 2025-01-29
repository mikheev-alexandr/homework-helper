import React from 'react';
import { useNavigate } from 'react-router-dom';
import '../styles/HomePage.css';

export default function HomePage() {
    const navigate = useNavigate();
  
    const handleTeacherSignIn = () => navigate('/auth/teacher/sign-in');
    const handleTeacherSignUp = () => navigate('/auth/teacher/sign-up');
    const handleStudentSignIn = () => navigate('/auth/student/sign-in');
  
    return (
      <div className="homepage">
        <header className="homepage-header">
          <h1>Добро пожаловать в систему проверки домашних заданий</h1>
          <p>Удобный сервис для управления образовательным процессом.</p>
        </header>
  
        <main className="homepage-main">
          <section className="homepage-info">
            <h2>Для учителей</h2>
            <div className="info">
              <p>С помощью нашего сервиса вы можете:</p>
              <ul>
                <li>Создавать задания для своих учеников.</li>
                <li>Отслеживать прогресс выполнения домашних заданий.</li>
                <li>Оставлять отзывы и выставлять оценки.</li>
              </ul>
            </div>
            <div className="homepage-buttons">
              <button className="homepage-button" onClick={handleTeacherSignIn}>
                Войти
              </button>
              <button className="homepage-button" onClick={handleTeacherSignUp}>
                Зарегистрироваться
              </button>
            </div>
          </section>
  
          <section className="homepage-info">
            <h2>Для учеников</h2>
            <div className="info">
              <p>Наш сервис предоставляет ученикам возможность:</p>
              <ul>
                <li>Получать задания от учителей и загружать выполненные работы.</li>
                <li>Следить за своим прогрессом и оценками.</li>
                <li>Узнавать об ошибках и получать рекомендации.</li>
              </ul>
            </div>
            <div className="homepage-buttons">
              <button className="homepage-button" onClick={handleStudentSignIn}>
                Войти
              </button>
            </div>
          </section>
        </main>
  
        <footer className="homepage-footer">
          <p>© 2025 Система проверки домашних заданий</p>
        </footer>
      </div>
    );
}