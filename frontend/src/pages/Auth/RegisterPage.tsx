import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Input } from '../../components/UI/Input';
import { Button } from '../../components/UI/Button';
import { useAuth } from '../../contexts/useAuth';
import { useNotifications } from '../../contexts/useNotifications';
import { getGoogleAuthUrlApi } from '../../api/auth';
import styles from './RegisterPage.module.css';

export default function RegisterPage() {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [isLoading, setIsLoading] = useState(false);
    const [googleLoading, setGoogleLoading] = useState(false);
    const { register } = useAuth();
    const navigate = useNavigate();
    const { addNotification } = useNotifications();

    const validate = (): boolean => {
        const newErrors: Record<string, string> = {};
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (name.length == 0) newErrors.name = 'Введите имя пользователя';
        if (!emailRegex.test(email)) newErrors.email = 'Введите корректный email';
        if (password.length < 8) newErrors.password = 'Пароль должен быть не менее 8 символов';
        if (password !== confirmPassword) newErrors.confirmPassword = 'Пароли не совпадают';

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

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;

        setIsLoading(true);
        try {
            await register({ name, email, password });
            navigate('/login');
        } catch {
            addNotification('error', 'Ошибка регистрации. Попробуйте позже.');
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <div className={styles.form}>
            <h1 className={styles.title}>Регистрация в NexaFlow</h1>

            <button className={styles.googleBtn} onClick={handleGoogleLogin} disabled={googleLoading}>
                <svg width="18" height="18" viewBox="0 0 48 48"><path fill="#EA4335" d="M24 9.5c3.54 0 6.71 1.22 9.21 3.6l6.85-6.85C35.9 2.38 30.47 0 24 0 14.62 0 6.51 5.38 2.56 13.22l7.98 6.19C12.43 13.72 17.74 9.5 24 9.5z"/><path fill="#4285F4" d="M46.98 24.55c0-1.57-.15-3.09-.38-4.55H24v9.02h12.94c-.58 2.96-2.26 5.48-4.78 7.18l7.73 6c4.51-4.18 7.09-10.36 7.09-17.65z"/><path fill="#FBBC05" d="M10.53 28.59A14.5 14.5 0 019.5 24c0-1.59.28-3.14.76-4.59l-7.98-6.19A23.99 23.99 0 000 24c0 3.77.87 7.35 2.56 10.56l7.97-5.97z"/><path fill="#34A853" d="M24 48c6.48 0 11.93-2.13 15.89-5.81l-7.73-6c-2.15 1.45-4.92 2.3-8.16 2.3-6.26 0-11.57-4.22-13.47-9.91l-7.98 5.97C6.51 42.62 14.62 48 24 48z"/></svg>
                {googleLoading ? 'Загрузка...' : 'Войти через Google'}
            </button>

            <div className={styles.divider}>или</div>

            <form className={styles.formInner} onSubmit={handleSubmit} noValidate>

            <Input
                label="Имя пользователя" name="name" value={name}
                onChange={(e) => setName(e.target.value)} error={errors.name}
                placeholder="username"
            />
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
            <Input
                label="Подтвердите пароль" name="confirmPassword" type="password" value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)} error={errors.confirmPassword}
                placeholder="••••••••"
            />

            <Button type="submit" variant="primary" disabled={isLoading}>
                {isLoading ? 'Регистрация...' : 'Зарегистрироваться'}
            </Button>

            <div className={styles.links}>
                <Link to="/login">Уже есть аккаунт? Войти</Link>
            </div>
        </form>
        </div>
    );
}
