import React, { useEffect, useState } from 'react';
import { useSearchParams } from 'react-router-dom';
import { confirmEmail } from '../../services/api.ts';
import '../../styles/ConfirmEmail.css';

export default function ConfirmEmail() {
  const [status, setStatus] = useState<'loading' | 'success' | 'error'>('loading');
  const [searchParams] = useSearchParams();

  useEffect(() => {
    const token = searchParams.get('token');
    if (token) {
      confirmEmail(token)
        .then(() => setStatus('success'))
        .catch(() => setStatus('error'));
    } else {
      setStatus('error');
    }
  }, [searchParams]);

  return (
    <div className="confirm-email-container">
      {status === 'loading' && <p>Processing your request...</p>}
      {status === 'success' && <p className="success">Ваш аккаунт зарегестрирован!</p>}
      {status === 'error' && <p className="error">Invalid or expired confirmation link.</p>}
    </div>
  );
}
