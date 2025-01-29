import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { getAssignments, deleteAssignment } from '../../services/api.ts';
import '../../styles/AssignmentsList.css';

interface Assignment {
  id: number;
  title: string;
  description: string;
  created_at: string;
}

export default function AssignmentsList() {
  const [assignments, setAssignments] = useState<Assignment[]>([]);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchAssignments() {
      try {
        const response = await getAssignments();
        setAssignments(response.data);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/teacher/sign-in');
        } else {
          console.error('Ошибка загрузки списка заданий:', error);
        }
      }
    }
    fetchAssignments();
  }, [navigate]);

  const handleDelete = async (id: number) => {
    try {
      await deleteAssignment(id);
      setAssignments(assignments.filter((assignment) => assignment.id !== id));
    } catch (error) {
      console.error('Ошибка удаления задания:', error);
    }
  };

  return (
    <div className="assignments-list-container">
      <h1>Список заданий</h1>
      <div className="actions-header">
        <Link to="/teacher/assignments/create" className="create-button">
          Создать задание
        </Link>
      </div>
      <ul>
        {assignments.map((assignment) => (
          <li key={assignment.id} className="assignment-item">
            <div className="assignment-info">
              <p className="title">{assignment.title}</p>
              <p className="description">{assignment.description}</p>
            </div>
            <div className="assignment-actions">
              <Link
                to={`/teacher/assignments/${assignment.id}`}
                className="details-button"
              >
                Подробнее
              </Link>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}
