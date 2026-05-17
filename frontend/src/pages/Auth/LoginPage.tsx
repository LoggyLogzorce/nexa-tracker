import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { Input } from '../../components/UI/Input';
import { Button } from '../../components/UI/Button';
import { useAuth } from '../../contexts/useAuth';
import { useNotifications } from '../../contexts/useNotifications';
import styles from './LoginPage.module.css';

export default function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [errors, setErrors] = useState<{ email?: string; password?: string }>({});
    const [isLoading, setIsLoading] = useState(false);
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
        <form className={styles.form} onSubmit={handleSubmit} noValidate>
            <h1 className={styles.title}>Вход в NexaFlow</h1>

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
    );
}
