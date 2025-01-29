import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { getHomeworksStudent } from '../../services/api.ts';

type Homework = {
  id: number;
  title: string;
  description: string;
  dueDate: string;
  status: string;
}

export default function Homeworks() {
    const [searchParams] = useSearchParams();
    const teacherId = searchParams.get('id');
    const [homeworks, setHomeworks] = useState<Homework[]>([]);
    const [teacherName, setTeacherName] = useState<string>(''); 
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchHomeworks() {
            try {
                const response = await getHomeworksStudent(teacherId ? parseInt(teacherId) : undefined);
                const data = response.data;
                const formattedData = data.map((item: any) => ({
                    id: item.id,
                    name: item.name,
                    title: item.title,
                    description: item.description,
                    dueDate: item.deadline,
                    status: item.status,
                }));
                setHomeworks(formattedData);

                if (formattedData.length > 0) {
                    setTeacherName(formattedData[0].name);
                }
            } catch (error: any) {
                if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                    navigate('/auth/student/sign-in');
                } else {
                    console.error('Ошибка загрузки домашних заданий:', error);
                }
            }
        }
        fetchHomeworks();
    }, [teacherId, navigate]);

  return (
          <div className="homework-list-container">
              <h1>{teacherId ? `Домашние задания - ${teacherName}` : 'Все домашние задания'}</h1>
              {homeworks.length === 0 ? (
                  <p>Нет заданных домашних заданий.</p>
              ) : (
                  <ul>
                      {homeworks.map(homework => (
                          <li key={homework.id} className="homework-item">
                              <h2>
                                  <Link to={`${homework.id}`} className="homework-link">
                                      {homework.title}
                                  </Link>
                              </h2>
                              <p>{homework.description}</p>
                              <p>Дата сдачи: {new Date(homework.dueDate).toLocaleDateString()}</p>
                              <p>Статус: {homework.status}</p>
                          </li>
                      ))}
                  </ul>
              )}
          </div>
      );
};