import React from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import '../styles/Header.css'
import { signOut } from '../services/api.ts';

export default function Header() {
  const location = useLocation();
  const navigate = useNavigate();

  const handleSignOut = async () => {
    await signOut();
    navigate('/');
  };

  const renderNavLinks = () => {
    if (location.pathname.startsWith('/teacher')) {
      return (
        <>
          <Link to="/teacher/students">Мои студенты</Link>
          <Link to="/teacher/homeworks">Домашние задания</Link>
          <Link to="/teacher/assignments">Задания</Link>
          <button onClick={handleSignOut}>Выход</button>
        </>
      );
    }

    if (location.pathname.startsWith('/student')) {
      return (
        <>
          <Link to="/student/teachers">Мои учителя</Link>
          <Link to="/student/homeworks">Мои задания</Link>
          <button onClick={handleSignOut}>Выход</button>
        </>
      );
    }

    if (location.pathname.startsWith('/auth')) {
      return <Link to="/">На главную</Link>;
    }

    return (
      <>
        <Link to="/">Главная</Link>
        <Link to="/auth/teacher/sign-in">Учитель</Link>
        <Link to="/auth/student/sign-in">Ученик</Link>
      </>
    );
  };

  return (
    <header className="app-header">
      <div className="container">
        <h1 className="logo">🏫 Homework Helper</h1>
        <nav className="nav-links">{renderNavLinks()}</nav>
      </div>
    </header>
  );
};