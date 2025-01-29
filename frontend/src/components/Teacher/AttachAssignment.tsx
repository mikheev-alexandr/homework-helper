import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../../styles/AttachAssignment.css';
import { attachHomework, getAssignments, getStudents } from '../../services/api.ts';

export default function AttachAssignment() {
  const [students, setStudents] = useState([]);
  const [assignments, setAssignments] = useState([]);
  const [selectedStudentId, setSelectedStudentId] = useState('');
  const [selectedAssignmentId, setSelectedAssignmentId] = useState('');
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [deadline, setDeadline] = useState('');
  const [message, setMessage] = useState(null);
  const [isSuccess, setIsSuccess] = useState(null);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchData() {
      try {
        const [studentsResponse, assignmentsResponse] = await Promise.all([
          getStudents(),
          getAssignments()
        ]);
        setStudents(studentsResponse.data);
        setAssignments(assignmentsResponse.data);
      } catch (error) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/teacher/sign-in');
        } else {
          console.error('Ошибка загрузки списка студентов:', error);
        }
      }
    }
    fetchData();
  }, []);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const payload = {
      assignment_id: parseInt(selectedAssignmentId, 10),
      student_id: parseInt(selectedStudentId, 10),
      title,
      description,
      deadline
    };

    try {
      await attachHomework(payload);
      setMessage('Задание успешно прикреплено!');
      setIsSuccess(true);
      setTimeout(() => navigate('/teacher/homeworks'), 2000);
    } catch (error) {
      console.error('Ошибка при прикреплении задания:', error);
      setMessage('Не удалось прикрепить задание. Попробуйте снова.');
      setIsSuccess(false);
    }
  };

  return (
    <div className="attach-homework-form">
      <h1>Прикрепить задание к ученику</h1>
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label htmlFor="student-select">Выберите ученика:</label>
          <select
            id="student-select"
            value={selectedStudentId}
            onChange={(e) => setSelectedStudentId(e.target.value)}
            required
          >
            <option value="">-- Выберите ученика --</option>
            {students.map((student) => (
              <option key={student.id} value={student.id}>
                {student.name}
              </option>
            ))}
          </select>
        </div>
  
        <div className="form-group">
          <label htmlFor="assignment-select">Выберите задание:</label>
          <select
            id="assignment-select"
            value={selectedAssignmentId}
            onChange={(e) => setSelectedAssignmentId(e.target.value)}
            required
          >
            <option value="">-- Выберите задание --</option>
            {assignments.map((assignment) => (
              <option key={assignment.id} value={assignment.id}>
                {assignment.title}
              </option>
            ))}
          </select>
        </div>
  
        <div className="form-group">
          <label htmlFor="title-input">Заголовок:</label>
          <input
            id="title-input"
            type="text"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
          />
        </div>
  
        <div className="form-group">
          <label htmlFor="description-textarea">Описание:</label>
          <textarea
            id="description-textarea"
            value={description}
            onChange={(e) => setDescription(e.target.value)}
          ></textarea>
        </div>
  
        <div className="form-group">
          <label htmlFor="deadline-input">Дедлайн:</label>
          <input
            id="deadline-input"
            type="datetime-local"
            value={deadline}
            onChange={(e) => setDeadline(e.target.value)}
            required
          />
        </div>
  
        <button type="submit">Прикрепить задание</button>
      </form>
  
      {message && (
        <p className={isSuccess ? 'success-message' : 'error-message'}>
          {message}
        </p>
      )}
    </div>
  );  
}
