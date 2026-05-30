import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { refresh } from '../../api/auth';
import { setCurrentAccessToken } from '../../api/client';
import { useAuth } from '../../contexts/useAuth';
import { useNotifications } from '../../contexts/useNotifications';

export default function GoogleCallbackPage() {
    const navigate = useNavigate();
    const { loadUser } = useAuth();
    const { addNotification } = useNotifications();

    useEffect(() => {
        const init = async () => {
            try {
                const { access_token } = await refresh();
                setCurrentAccessToken(access_token);
                await loadUser();
                navigate('/dashboard', { replace: true });
            } catch {
                addNotification('error', 'Ошибка при входе через Google');
                navigate('/login', { replace: true });
            }
        };
        init();
    }, [addNotification, loadUser, navigate]);

    return (
        <div style={{
            minHeight: '100vh',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            background: '#f3f4f6',
            color: '#6b7280',
            fontSize: '1rem',
        }}>
            Вход через Google...
        </div>
    );
}
