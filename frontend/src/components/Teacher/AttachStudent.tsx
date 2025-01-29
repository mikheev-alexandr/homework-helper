import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { attachStudent } from '../../services/api.ts';
import '../../styles/AttachStudent.css'

export default function AttachStudent() {
  const [action, setAction] = useState('create');
  const [name, setName] = useState('');
  const [classNumber, setClassNumber] = useState(1);
  const [codeWord, setCodeWord] = useState('');
  const [responseMessage, setResponseMessage] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    const input = {
      action,
      name: action === 'create' ? name : undefined,
      class_number: action === 'create' ? classNumber : undefined,
      code_word: action === 'attach' ? codeWord : undefined,
    };

    try {
      const response = await attachStudent(input);
      if (action === 'create') {
        setResponseMessage(
          `Успешно создан! Логин: ${response.data.login} Пароль: ${response.data.password}`
        );
      } else if (action === 'attach') {
        setResponseMessage(
          `Студент присоединён! Имя: ${response.data.name}, Кодовое слово: ${response.data.code_word}, Класс: ${response.data.class_number}`
        );
      }
    } catch (error: any) {
      if (error.response && (error.response.status === 401 || error.response.status === 403)) {
        navigate('/auth/teacher/sign-in');
      } else {
        console.error('Ошибка загрузки списка студентов:', error);
      }
    }
    
  };

  return (
    <div className="attach-student-container">
      <h1>Добавление студента</h1>
      <form onSubmit={handleSubmit}>
        <div>
          <label>
            Выберите действие:
            <select
              value={action}
              onChange={(e) => setAction(e.target.value)}
            >
              <option value="create">Создать нового студента</option>
              <option value="attach">Присоединить существующего студента</option>
            </select>
          </label>
        </div>

        {action === 'create' && (
          <>
            <div>
              <label>
                Имя:
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  required
                />
              </label>
            </div>
            <div>
              <label>
                Класс:
                <input
                  type="number"
                  value={classNumber}
                  onChange={(e) => setClassNumber(Number(e.target.value))}
                  min="1"
                  required
                />
              </label>
            </div>
          </>
        )}

        {action === 'attach' && (
          <div>
            <label>
              Кодовое слово:
              <input
                type="text"
                value={codeWord}
                onChange={(e) => setCodeWord(e.target.value)}
                required
              />
            </label>
          </div>
        )}

        <button type="submit">Подтвердить</button>
      </form>

      {responseMessage && (
        <div className="response-message">
          <p>{responseMessage}</p>
        </div>
      )}

      <button onClick={() => navigate('/teacher/students')} className="back-button">
        Вернуться к списку студентов
      </button>
    </div>
  );
}
