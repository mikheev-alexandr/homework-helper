import React, { useEffect, useState } from 'react';
import { useSearchParams, useNavigate, Link } from 'react-router-dom';
import { getHomeworksTeacher } from '../../services/api.ts';
import '../../styles/HomeworkList.css';

type Homework = {
    id: number;
    name: string;
    class_number: number;
    title: string;
    description: string;
    dueDate: string;
    status: string;
};

export default function HomeworkTeacherList() {
    const [searchParams] = useSearchParams();
    const studentId = searchParams.get('id');
    const [homeworks, setHomeworks] = useState<Homework[]>([]);
    const [studentName, setStudentName] = useState<string>(''); 
    const navigate = useNavigate();

    useEffect(() => {
        async function fetchHomeworks() {
            try {
                const response = await getHomeworksTeacher(studentId ? parseInt(studentId) : undefined);
                const data = response.data;
                const formattedData = data.map((item: any) => ({
                    id: item.id,
                    name: item.name,
                    class_number: item.class_number,
                    title: item.title,
                    description: item.description,
                    dueDate: item.deadline,
                    status: item.status,
                }));
                setHomeworks(formattedData);

                if (formattedData.length > 0) {
                    setStudentName(formattedData[0].name);
                }
            } catch (error: any) {
                if (error.response && (error.response.status === 401 || error.response.status === 403)) {
                    navigate('/auth/teacher/sign-in');
                } else {
                    console.error('Ошибка загрузки домашних заданий:', error);
                }
            }
        }
        fetchHomeworks();
    }, [studentId, navigate]);

    return (
        <div className="homework-list-container">
            <h1>{studentId ? `Домашние задания - ${studentName}` : 'Все домашние задания'}</h1>
            <div className="actions-header">
                <Link to={`/teacher/homeworks/attach${studentId ? `?id=${studentId}` : ''}`} className="create-homework-button">
                    Задать новое задание
                </Link>
            </div>
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
                            <p>{homework.name} Класс: {homework.class_number}</p>
                            <p>{homework.description}</p>
                            <p>Дата сдачи: {new Date(homework.dueDate).toLocaleDateString()}</p>
                            <p>Статус: {homework.status}</p>
                        </li>
                    ))}
                </ul>
            )}
        </div>
    );
}
