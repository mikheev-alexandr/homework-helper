import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { getTeachers } from '../../services/api.ts';
import '../../styles/TeacherList.css';

interface Teacher {
  id: number;
  name: string;
}

export default function TeachersList() {
  const [teachers, setteachers] = useState<Teacher[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchteachers() {
      try {
        const response = await getTeachers();
        setteachers(response.data);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/student/sign-in');
        } else {
          console.error('Ошибка загрузки списка студентов:', error);
        }
      }
    }
    fetchteachers();
  }, [navigate]);

  return (
    <div className="teachers-list-container">
      <h1>Список учителей</h1>
      <ul>
        {teachers.map(teacher => (
          <li key={teacher.id} className="teacher-item">
            <div className="teacher-info">
              <p className="name">{teacher.name}</p>
            </div>
            <div className="teacher-actions">
              <Link to={`/student/homeworks?id=${teacher.id}`} className="details-button">Домашние задания</Link>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
