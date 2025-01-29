import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import {
  getStudentHomeworkById,
  attachStudentHomework,
  deleteStudentHomework,
  updateStudentHomework,
} from '../../services/api.ts';
import { Link } from 'react-router-dom';
import '../../styles/HomeworkDetails.css';

export default function HomeworkStudentDetails() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [homework, setHomework] = useState<any>(null);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState<FormData>(new FormData());
  const [formDescription, setFormDescription] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [files, setFiles] = useState<File[]>([]);
  const [loading, setLoading] = useState(true);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);

  useEffect(() => {
    async function fetchHomework() {
      try {
        const response = await getStudentHomeworkById(Number(id));
        const data = response.data;
        setHomework(data);
        setLoading(false);
        setError(false);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/student/sign-in');
        } else {
          setError('Ошибка загрузки задания.');
        }
      }
    }
    fetchHomework();
  }, [id]);

  const handleAttach = async () => {
    try {
      const formData = new FormData();
      formData.append('text', formDescription);
  
      files.forEach((file) => formData.append('files', file));
  
      await attachStudentHomework(Number(id), formData);
      const response = await getStudentHomeworkById(Number(id));
      setHomework(response.data);
      setIsSubmitting(false);
    } catch {
      setError('Не удалось прикрепить решение.');
    }
  };
  

  const handleUpdate = async () => {
    try {
      const formData = new FormData();
      formData.append('text', formDescription);
  
      if (files.length > 0) {
        files.forEach((file) => formData.append('files', file));
      } else {
        formData.append('files', new Blob([]), '');
      }

      const response = await updateStudentHomework(Number(homework.sub_id), formData);

      if (response.data.updated) {
        const response = await getStudentHomeworkById(Number(id))
        setHomework(response.data)
        setIsEditing(false);
      } else {
        setError('Не удалось обновить задание. Заданние уже оценено.');
        const response = await getStudentHomeworkById(Number(id))
        setHomework(response.data)
        setIsEditing(false);
      }
    } catch (err) {
      setError('Ошибка при обновлении задания.');
    }
  };

  const handleDelete = async () => {
    try {
      const response = await deleteStudentHomework(Number(homework.sub_id));
      if (response.data.deleted) {
        const response = await getStudentHomeworkById(Number(id))
        setHomework(response.data)
      } else {
        setError('Не удалось удалить задание. Задание уже оценено.');
        const response = await getStudentHomeworkById(Number(id))
        setHomework(response.data)
        setIsEditing(false);
      }
    } catch (err) {
      setError('Ошибка при удалении задания.');
    }
  };

  if (loading) return (
    <div className="not-found-container">
      <h1>404</h1>
      <p>Упс! Страница не найдена.</p>
      <p>
        Возможно, вы ошиблись в адресе или страница была удалена. Вы можете вернуться на{' '}
        <Link to="/student/homeworks" className="link">к списку домашних заданий</Link>.
      </p>
    </div>
  );


  const homeworkFiles = homework?.hw_files || [];
  const submissionFiles = homework?.sub_files || [];

  return (
    <div className="teacher-homework-details-container">
      {error && <div className="error-block">{error}</div>}

      {isSubmitting ? (
        <div>
          <h1>Прикрепить решение</h1>
          <textarea
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder="Комментарий"
          />
          <input
            type="file"
            multiple
            onChange={(e) => {
              const selectedFiles = e.target.files;
              if (selectedFiles) {
                setFiles(Array.from(selectedFiles));
              }
            }}
          />
          <button onClick={handleAttach}>Сохранить</button>
          <button onClick={() => setIsSubmitting(false)}>Отменить</button>
        </div>
      ) : isEditing ? (
        <div>
          <h1>Редактировать решение</h1>
          <textarea
            value={formDescription || homework?.text || ''}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder="Комментарий"
          />
          <input
            type="file"
            multiple
            onChange={(e) => {
              const selectedFiles = e.target.files;
              if (selectedFiles) {
                setFiles(Array.from(selectedFiles));
              }
            }}
          />
          <button onClick={handleUpdate}>Сохранить изменения</button>
          <button onClick={() => setIsEditing(false)}>Отменить</button>
        </div>
      ) : (
        <div>
          <h1>{homework.title}</h1>
          <p>{homework.description}</p>
          <p>Сделать до: {new Date(homework.deadline).toLocaleString()}</p>
          <p>Статус: {homework.status}</p>
          <ul>
            {homeworkFiles.length > 0 ? (
              homeworkFiles.map((file, index) => (
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

          {homework.status === 'не решено' ? (
            <button onClick={() => setIsSubmitting(true)}>Прикрепить решение</button>
          ) : (
            <div className="solution-details">
              <h2>Решение</h2>
              <p>Комментарий: {homework.text}</p>
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
              
              {homework.status == 'оценено' && (
                <div className="grade-details">
                  <h2>Оценка: {homework.grade}</h2>
                  <p><strong>Комментарий:</strong> {homework.feedback}</p>
                </div>
              )}

              <button onClick={() => setIsEditing(true)}>Редактировать решение</button>
              <button onClick={() => setIsDeleteModalOpen(true)}>Удалить решение</button>
            </div>
          )}
        </div>
      )}

      {isDeleteModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Подтверждение удаления</h2>
            <p>Вы уверены, что хотите удалить это решение?</p>
            <button className="confirm-button" onClick={handleDelete}>
              Да
            </button>
            <button className="cancel-button" onClick={() => setIsDeleteModalOpen(false)}>
              Отмена
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
