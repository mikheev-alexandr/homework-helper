import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { deleteHomework, getTeacherHomeworkById, gradeHomework, updateHomework } from '../../services/api.ts';
import { Link } from 'react-router-dom';

export default function HomeworkTeacherDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [homework, setHomework] = useState<any>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [isGrading, setIsGrading] = useState(false);
  const [formData, setFormData] = useState({
    title: '',
    description: '',
    deadline: '',
    files: null as FileList | null,
  });
  const [gradeData, setGradeData] = useState({ grade: 0, feedback: '' });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  useEffect(() => {
    async function fetchHomework() {
      try {
        const response = await getTeacherHomeworkById(Number(id));
        const data = response.data;
        setHomework(data);
        setFormData({
          title: data.title,
          description: data.description,
          deadline: data.deadline,
          files: null,
        });
        setLoading(false);
        setError(false);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/teacher/sign-in');
        } else {
          console.error('Ошибка загрузки задания:', error);
        }
      }
    }
    fetchHomework();
  }, [id]);

  const refreshHomework = async () => {
    try {
      const response = await getTeacherHomeworkById(Number(id));
      setHomework(response.data);
    } catch (error) {
      console.error('Ошибка обновления данных:', error);
    }
  };

  const handleUpdate = async () => {
    try {
      const formDataToSend = new FormData();
      formDataToSend.append('title', formData.title);
      formDataToSend.append('description', formData.description);
      formDataToSend.append('deadline', formData.deadline);
  
      const response = await updateHomework(Number(id), formDataToSend);
  
      if (response.data.updated) {
        await refreshHomework();
        setIsEditing(false);
      } else {
        setError('Не удалось обновить задание. Заданние уже решено.');
        await refreshHomework();
        setIsEditing(false);
      }
    } catch (err) {
      setError('Ошибка при обновлении задания.');
    }
  };
  

  const handleDelete = async () => {
    try {
      const response = await deleteHomework(Number(id));
      if (response.data.deleted) {
        await refreshHomework();
        setIsEditing(false);
      } else {
        setError('Не удалось удалить задание. Задание уже решено.');
        await refreshHomework();
        setIsEditing(false);
      }
    } catch (err) {
      setError('Ошибка при обновлении задания.');
    }
  };

  const handleGradeSubmit = async () => {
    try {
      await gradeHomework(Number(id), gradeData);
      await refreshHomework();
      setIsGrading(false);
      setGradeData({ grade: 0, feedback: '' });
    } catch (err) {
      setError('Не удалось оценить задание.');
    }
  };

  if (loading) return (
    <div className="not-found-container">
      <h1>404</h1>
      <p>Упс! Страница не найдена.</p>
      <p>
        Возможно, вы ошиблись в адресе или страница была удалена. Вы можете вернуться на{' '}
        <Link to="/teacher/homeworks" className="link">к списку домашних заданий</Link>.
      </p>
    </div>
  );

  const files = homework?.hw_files || [];
  const submissionFiles = homework?.sub_files || [];

  return (
    <div className="teacher-homework-details-container">

      {error && <div className="error-block">{error}</div>}

      {isEditing ? (
        <div>
          <h1>Редактировать задание</h1>
          <input
            type="text"
            value={formData.title}
            onChange={(e) => setFormData({ ...formData, title: e.target.value })}
            placeholder="Название"
          />
          <textarea
            value={formData.description}
            onChange={(e) => setFormData({ ...formData, description: e.target.value })}
            placeholder="Описание"
          />
          <input
            type="datetime-local"
            value={formData.deadline}
            onChange={(e) => setFormData({ ...formData, deadline: e.target.value })}
          />
          <button onClick={handleUpdate}>Сохранить изменения</button>
          <button onClick={() => setIsEditing(false)}>Отменить</button>
        </div>
      ) : isGrading ? (
        <div className="grading-container">
          <h2>Оценить задание</h2>
          <input
            type="number"
            value={gradeData.grade}
            onChange={(e) => setGradeData({ ...gradeData, grade: parseInt(e.target.value) })}
            placeholder="Оценка (1-5)"
          />
          <textarea
            value={gradeData.feedback}
            onChange={(e) => setGradeData({ ...gradeData, feedback: e.target.value })}
            placeholder="Обратная связь"
          />
          <button onClick={handleGradeSubmit}>Отправить оценку</button>
          <button onClick={() => setIsGrading(false)}>Отменить</button>
        </div>
      ) : (
        <div>
          <h1>{homework.title}</h1>
          <p>{homework.description}</p>
          <p>Сделать до: {new Date(homework.deadline).toLocaleString()}</p>
          <p>Статус: {homework.status}</p>
          <ul>
            {files.length > 0 ? (
              files.map((file, index) => (
                <li key={index}>
                  <a
                    className="file-link"
                    href={`data:application/octet-stream;base64,${file.content}`}
                    download={file.filePath}
                  >
                    {file.filePath}
                  </a>
                </li>
              ))
            ) : (
              <p className="no-files-message">Нет файлов для загрузки.</p>
            )}
          </ul>

          {(homework.status === 'решено' || homework.status === 'оценено') && (
            <div className="solution-details">
              <h2>Решение</h2>
              <p>Текст решения: {homework.text}</p>
              <p>Отправлено: {new Date(homework.submited_at).toLocaleString()}</p>
              <ul>
                {submissionFiles.length > 0 ? (
                  submissionFiles.map((file, index) => (
                    <li key={index}>
                      <a
                        className="file-link"
                        href={`data:application/octet-stream;base64,${file.content}`}
                        download={file.filePath}
                      >
                        {file.filePath}
                      </a>
                    </li>
                  ))
                ) : (
                  <p className="no-files-message">Нет файлов для решения.</p>
                )}
              </ul>
            </div>
          )}

          {homework.status === 'оценено' && (
            <div className="grade-details">
              <h2>Оценка: {homework.grade}</h2>
              <p><strong>Комментарий:</strong> {homework.feedback}</p>
            </div>
          )}

          <button onClick={() => setIsEditing(true)}>Изменить задание</button>
          <button onClick={() => setIsGrading(true)}>Оценить задание</button>
          <button onClick={() => setIsDeleteModalOpen(true)}>Удалить задание</button>
        </div>
      )}

  
      {isDeleteModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Подтверждение удаления</h2>
            <p>Вы уверены, что хотите удалить это задание?</p>
            <button className="confirm-button" onClick={handleDelete}>Да</button>
            <button className="cancel-button" onClick={() => setIsDeleteModalOpen(false)}>Отмена</button>
          </div>
        </div>
      )}
    </div>
  );  
}
