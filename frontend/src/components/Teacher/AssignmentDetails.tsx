import React, { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getAssignment, updateAssignment, deleteAssignment } from '../../services/api.ts';
import '../../styles/AssignmentDetails.css';

interface Assignment {
  id: number;
  title: string;
  description: string;
  created_at: string;
  files: { filePath: string; content: string }[];
}

export default function AssignmentDetails() {
  const { id } = useParams<{ id: string }>();
  const [assignment, setAssignment] = useState<Assignment | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [formTitle, setFormTitle] = useState('');
  const [formDescription, setFormDescription] = useState('');
  const [newFiles, setNewFiles] = useState<File[]>([]);
  const [message, setMessage] = useState<string | null>(null);
  const [isSuccess, setIsSuccess] = useState<boolean | null>(null);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    async function fetchAssignment() {
      try {
        const response = await getAssignment(Number(id));
        setAssignment(response.data);
        setFormTitle(response.data.title);
        setFormDescription(response.data.description);
      } catch (error: any) {
        if (error.response && (error.response.status === 401 || error.response.status === 403)) {
          navigate('/auth/teacher/sign-in');
        } else if (error.response && (error.response.status === 500 && error.response.data.error === 'sql: no rows in result set')) {
          navigate('/teacher/assignments')
        } else {
          console.error('Ошибка загрузки задания:', error);
        }
      }
    }
    fetchAssignment();
  }, [navigate, id]);

  const handleEditToggle = () => {
    setIsEditing(!isEditing);
    setMessage(null);
    setIsSuccess(null);
  };

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files) {
      setNewFiles(Array.from(event.target.files));
    }
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();

    const formData = new FormData();
    formData.append('title', formTitle);
    formData.append('description', formDescription);
    newFiles.forEach((file) => formData.append('files', file));

    try {
      await updateAssignment(Number(id), formData);
      setMessage('Задание успешно обновлено!');
      setIsSuccess(true);
      setIsEditing(false);

      const response = await getAssignment(Number(id));
      setAssignment(response.data);
    } catch (error) {
      console.error('Ошибка при обновлении задания:', error);
      setMessage('Не удалось обновить задание. Попробуйте снова.');
      setIsSuccess(false);
    }
  };

  const handleDelete = async () => {
    if (!id) return;

    try {
      await deleteAssignment(Number(id));
      setMessage('Задание успешно удалено!');
      setIsSuccess(true);
      navigate('/teacher/assignments');
    } catch (error) {
      console.error('Ошибка при удалении задания:', error);
      setMessage('Не удалось удалить задание. Попробуйте снова.');
      setIsSuccess(false);
    }
  };

  const files = assignment?.files || [];

  return (
    <div className="create-assignment-form">
      {isEditing ? (
        <form onSubmit={handleSubmit}>
          <h1>Редактировать задание</h1>
          <div>
            <label>
              Название:
              <input
                type="text"
                value={formTitle}
                onChange={(e) => setFormTitle(e.target.value)}
                required
              />
            </label>
          </div>
          <div>
            <label>
              Описание:
              <textarea
                value={formDescription}
                onChange={(e) => setFormDescription(e.target.value)}
                required
              />
            </label>
          </div>
          <div>
            <label>
              Добавить файлы:
              <input type="file" multiple onChange={handleFileChange} />
            </label>
          </div>
          <button type="submit">Сохранить изменения</button>
          <button type="button" onClick={handleEditToggle}>
            Отменить
          </button>
        </form>
      ) : (
        <>
          <h1>{assignment?.title}</h1>
          <p>{assignment?.description}</p>
          <p className="assignment-created-at">
            Создано: {new Date(assignment?.created_at ?? '').toLocaleString()}
          </p>
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
          <button className="edit-button" onClick={handleEditToggle}>
            Редактировать
          </button>
          <button
            className="delete-button"
            onClick={() => setIsDeleteModalOpen(true)}
          >
            Удалить
          </button>
        </>
      )}
      {message && (
        <p className={isSuccess ? 'success-message' : 'error-message'}>
          {message}
        </p>
      )}

      {isDeleteModalOpen && (
        <div className="modal">
          <div className="modal-content">
            <h2>Подтверждение удаления</h2>
            <p>Вы уверены, что хотите удалить это задание?</p>
            <button className="confirm-button" onClick={handleDelete}>
              Да
            </button>
            <button
              className="cancel-button"
              onClick={() => setIsDeleteModalOpen(false)}
            >
              Отмена
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
