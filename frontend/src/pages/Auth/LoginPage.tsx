import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Input } from '../../components/UI/Input';
import { Button } from '../../components/UI/Button';
import { useAuth } from '../../contexts/useAuth';
import { useNotifications } from '../../contexts/useNotifications';
import { getGoogleAuthUrlApi, getYandexAuthUrlApi } from '../../api/auth';
import styles from './LoginPage.module.css';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errors, setErrors] = useState<{ email?: string; password?: string }>({});
    const [isLoading, setIsLoading] = useState(false);
    const [googleLoading, setGoogleLoading] = useState(false);
    const [yandexLoading, setYandexLoading] = useState(false);
    const { login } = useAuth();
    const { addNotification } = useNotifications();

    const validate = (): boolean => {
        const newErrors: typeof errors = {};
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (!emailRegex.test(email)) newErrors.email = 'Введите корректный email';
        if (!password.trim()) newErrors.password = 'Пароль не может быть пустым';

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleGoogleLogin = async () => {
        setGoogleLoading(true);
        try {
            const url = await getGoogleAuthUrlApi();
            window.location.href = url;
        } catch {
            addNotification('error', 'Ошибка при входе через Google');
            setGoogleLoading(false);
        }
    };

    const handleYandexLogin = async () => {
        setYandexLoading(true);
        try {
            const url = await getYandexAuthUrlApi();
            window.location.href = url;
        } catch {
            addNotification('error', 'Ошибка при входе через Яндекс');
            setYandexLoading(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;

        setIsLoading(true);
        try {
            await login({ email, password });
        } catch {
            addNotification('error', 'Неверный email или пароль');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className={styles.form}>
            <h1 className={styles.title}>Вход в NexaFlow</h1>

            <div className={styles.socialRow}>
                <button className={styles.googleBtn} onClick={handleGoogleLogin} disabled={googleLoading}>
                    <svg width="18" height="18" viewBox="0 0 48 48"><path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/><path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/><path fill="#FBBC05" d="M10.53 28.59A14.5 14.5 0 019.5 24c0-1.59.28-3.14.76-4.59l-7.98-6.19A23.99 23.99 0 000 24c0 3.77.87 7.35 2.56 10.56l7.97-5.97z"/><path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 5.97C6.51 42.62 14.62 48 24 48z"/></svg>
                    {googleLoading ? 'Загрузка...' : 'Google'}
                </button>
                <button className={styles.yandexBtn} onClick={handleYandexLogin} disabled={yandexLoading}>
                    <svg fill="currentColor" fill-rule="evenodd" height="18" viewBox="0 0 24 24" width="18" xmlns="http://www.w3.org/2000/svg"><title>Yandex</title><path d="M16.376 12.644L21 2h-3.842l-4.624 10.644h3.842zM13.915 24v-3.733c0-2.822-.352-3.64-1.407-5.988L6.933 2H3l7.124 15.709V24h3.79z"></path></svg>
                    {yandexLoading ? 'Загрузка...' : 'Яндекс'}
                </button>
            </div>

            <div className={styles.divider}>или</div>

            <form className={styles.formInner} onSubmit={handleSubmit} noValidate>
            <Input
                label="Email" name="email" type="email" value={email}
                onChange={(e) => setEmail(e.target.value)} error={errors.email}
                placeholder="you@example.com"
            />

            <Input
                label="Пароль" name="password" type="password" value={password}
                onChange={(e) => setPassword(e.target.value)} error={errors.password}
                placeholder="••••••••"
            />

            <Button type="submit" variant="primary" disabled={isLoading}>
                {isLoading ? 'Вход...' : 'Войти'}
            </Button>

            <div className={styles.links}>
                <Link to="/forgot-password">Забыли пароль?</Link>
                <Link to="/register">Регистрация</Link>
            </div>
        </form>
        </div>
    );
}
