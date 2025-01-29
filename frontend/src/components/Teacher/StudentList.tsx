import React, { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { getStudents, deleteStudent } from '../../services/api.ts';
import '../../styles/StudentList.css';

interface Student {
  id: number;
  name: string;
  class_number: number;
  code_word: string;
}

export default function StudentsList() {
  const [students, setStudents] = useState<Student[]>([]);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [studentToDelete, setStudentToDelete] = useState<Student | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchStudents() {
      try {
        const response = await getStudents();
        setStudents(response.data);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/teacher/sign-in');
        } else {
          console.error('Ошибка загрузки списка студентов:', error);
        }
      }
    }
    fetchStudents();
  }, [navigate]);

  const openDeleteModal = (student: Student) => {
    setStudentToDelete(student);
    setIsDeleteModalOpen(true);
  };

  const closeDeleteModal = () => {
    setStudentToDelete(null);
    setIsDeleteModalOpen(false);
  };

  const handleDelete = async () => {
    if (!studentToDelete) return;

    try {
      await deleteStudent(studentToDelete.id);
      setStudents(students.filter(student => student.id !== studentToDelete.id));
      closeDeleteModal();
    } catch (error: any) {
      if (error.response && (error.response.status === 401 || error.response.status === 403)) {
        navigate('/auth/teacher/sign-in');
      } else {
        console.error('Ошибка удаления студента:', error);
      }
    }
  };

  return (
    <div className="students-list-container">
      <h1>Список студентов</h1>
      <Link to="/teacher/students/attach" className="button">Добавить студента</Link>
      <ul>
        {students.map(student => (
          <li key={student.id} className="student-item">
            <div className="student-info">
              <p className="name">{student.name}</p>
              <p className="class">Класс: {student.class_number}</p>
              <p className="code-word">Кодовое слово: {student.code_word}</p>
            </div>
            <div className="student-actions">
              <Link to={`/teacher/homeworks?id=${student.id}`} className="details-button">Домашние задания</Link>
              <button onClick={() => openDeleteModal(student)} className="delete-button">Удалить</button>
            </div>
          </li>
        ))}
      </ul>

      {isDeleteModalOpen && studentToDelete && (
        <div className="modal">
          <div className="modal-content">
            <h2>Подтверждение удаления</h2>
            <p>Вы уверены, что хотите удалить студента "{studentToDelete.name}"?</p>
            <button className="confirm-button" onClick={handleDelete}>Да</button>
            <button className="cancel-button" onClick={closeDeleteModal}>Отмена</button>
          </div>
        </div>
      )}
    </div>
  );
}
