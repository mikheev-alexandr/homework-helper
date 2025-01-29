import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { createAssignment } from '../../services/api.ts';
import '../../styles/CreateAssignment.css'

export default function CreateAssignment() {
  const [title, setTitle] = useState('');
  const [description, setDescription] = useState('');
  const [files, setFiles] = useState<FileList | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [isSuccess, setIsSuccess] = useState<boolean | null>(null);
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append('title', title);
    formData.append('description', description);

    if (files) {
      Array.from(files).forEach((file) => {
        formData.append('files', file as Blob);
      });
    }

    try {
      await createAssignment(formData);
      setMessage('Задание успешно создано!');
      setIsSuccess(true);
      setTitle('');
      setDescription('');
      setFiles(null);
    } catch (error: any) {
      if (error.response && (error.response.status === 401 || error.response.status === 403)) {
        navigate('/auth/teacher/sign-in');
      } else {
        console.error('Ошибка загрузки списка студентов:', error);
      }
    }
  };

  return (
    <div className="create-assignment-form">
      <form onSubmit={handleSubmit}>
        <h1>Создать задание</h1>
        <div>
          <label>
            Название:
            <input
              type="text"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              required
            />
          </label>
        </div>
        <div>
          <label>
            Описание:
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
            />
          </label>
        </div>
        <div>
          <label>
            Файлы:
            <input
              type="file"
              multiple
              onChange={(e) => setFiles(e.target.files)}
            />
          </label>
        </div>
        <button type="submit">Создать</button>
      </form>
      {message && (
        <p className={isSuccess ? 'success-message' : 'error-message'}>
          {message}
        </p>
      )}
    </div>
  );
};