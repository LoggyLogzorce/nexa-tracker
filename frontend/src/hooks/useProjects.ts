import { useState } from 'react';
import type { Project, ProjectStatus, Priority } from '../types/project';

const mockProjects: Project[] = [
    { id: '1', title: 'Редизайн мобильного приложения', description: 'Обновление UI/UX дизайна для iOS и Android платформ, включая создание дизайн-системы.', status: 'В работе', priority: 'Высокий', role: 'owner', createdAt: '2026-05-12' },
    { id: '2', title: 'Интеграция платежной системы', description: 'Подключение Stripe API и настройка webhook событий для обработки платежей.', status: 'Завершен', priority: 'Низкий', role: 'owner', createdAt: '2026-05-01' },
    { id: '3', title: 'Аналитика пользователей', description: 'Внедрение Google Analytics 4 и настройка кастомных событий для трекинга поведения.', status: 'Планирование', priority: 'Средний', role: 'member', createdAt: '2026-05-14' },
    { id: '4', title: 'Миграция базы данных', description: 'Перенос данных из PostgreSQL в MongoDB для улучшения производительности чтения.', status: 'В работе', priority: 'Высокий', role: 'owner', createdAt: '2026-05-10' },
    { id: '5', title: 'Разработка API Gateway', description: 'Создание единой точки входа для микросервисов с Rate Limiting и аутентификацией.', status: 'В работе', priority: 'Высокий', role: 'member', createdAt: '2026-05-08' },
    { id: '6', title: 'DevOps инфраструктура', description: 'Настройка CI/CD пайплайнов и мониторинга с использованием Kubernetes и Prometheus.', status: 'Завершен', priority: 'Средний', role: 'read_only', createdAt: '2026-05-05' },
];

export function useProjects() {
    const [projects, setProjects] = useState<Project[]>(mockProjects);

    const myProjects = projects.filter(p => p.role === 'owner');

    const createProject = (data: { title: string; description: string; status: ProjectStatus; priority: Priority }) => {
        const newProject: Project = {
            id: String(Date.now()),
            ...data,
            role: 'owner',
            createdAt: new Date().toISOString().split('T')[0],
        };
        setProjects(prev => [newProject, ...prev]);
    };

    const updateProject = (updated: Project) => {
        setProjects(prev => prev.map(p => p.id === updated.id ? updated : p));
    };

    const deleteProject = (id: string) => {
        setProjects(prev => prev.filter(p => p.id !== id));
    };

    return { projects, myProjects, createProject, updateProject, deleteProject };
}
