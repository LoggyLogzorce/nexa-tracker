import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Input } from '../../components/UI/Input';
import { Button } from '../../components/UI/Button';
import { useAuth } from '../../contexts/useAuth';
import styles from './RegisterPage.module.css';

export default function RegisterPage() {
    const [name, setName] = useState('');
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [confirmPassword, setConfirmPassword] = useState('');
    const [errors, setErrors] = useState<Record<string, string>>({});
    const [isLoading, setIsLoading] = useState(false);
    const { register } = useAuth();
    const navigate = useNavigate();

    const validate = (): boolean => {
        const newErrors: Record<string, string> = {};
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

        if (!emailRegex.test(email)) newErrors.email = 'Введите корректный email';
        if (password.length < 6) newErrors.password = 'Пароль должен быть не менее 6 символов';
        if (password !== confirmPassword) newErrors.confirmPassword = 'Пароли не совпадают';

        setErrors(newErrors);
        return Object.keys(newErrors).length === 0;
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!validate()) return;

        setIsLoading(true);
        try {
            await register({ name, email, password });
            navigate('/login');
        } catch {
            setErrors({ general: 'Ошибка регистрации. Попробуйте позже.' });
        } finally {
            setIsLoading(false);
        }
    };

    return (
        <form className={styles.form} onSubmit={handleSubmit} noValidate>
            <h1 className={styles.title}>Регистрация в NexaFlow</h1>

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
    );
}
