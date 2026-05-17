import { useState, useEffect, useCallback } from 'react';
import type { Project, ProjectStatus, Priority } from '../types/project';
import { getProjects, getOwnedProjects, createProjectApi, updateProjectApi, deleteProjectApi } from '../api/projects';
import { useNotifications } from '../contexts/useNotifications';

export function useProjects(mode: 'all' | 'owned' = 'all') {
    const { addNotification } = useNotifications();
    const [projects, setProjects] = useState<Project[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetch = mode === 'all' ? getProjects : getOwnedProjects;

        fetch()
            .then(data => { setProjects(data); setError(null); })
            .catch(err => setError(err instanceof Error ? err.message : 'Ошибка загрузки проектов'))
            .finally(() => setIsLoading(false));
    }, [mode]);

    const createProject = useCallback((data: { title: string; description: string; status: ProjectStatus; priority: Priority }) => {
        return createProjectApi(data)
            .then(project => { setProjects(prev => [project, ...prev]); addNotification('success', 'Проект успешно создан'); })
            .catch(err => { addNotification('error', 'Ошибка при создании проекта'); throw err; });
    }, [addNotification]);

    const updateProject = useCallback((updated: Project) => {
        updateProjectApi(updated.id, updated)
            .then(project => { setProjects(prev => prev.map(p => p.id === project.id ? project : p)); addNotification('success', 'Проект успешно обновлён'); })
            .catch(() => addNotification('error', 'Ошибка при обновлении проекта'));
    }, [addNotification]);

    const deleteProject = useCallback((id: string) => {
        deleteProjectApi(id)
            .then(() => { setProjects(prev => prev.filter(p => p.id !== id)); addNotification('success', 'Проект успешно удалён'); })
            .catch(() => addNotification('error', 'Ошибка при удалении проекта'));
    }, [addNotification]);

    return { projects, isLoading, error, createProject, updateProject, deleteProject };
}
