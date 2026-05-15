export type ProjectStatus = 'В работе' | 'Планирование' | 'Завершен';
export type Priority = 'Низкий' | 'Средний' | 'Высокий';
export type ProjectRole = 'owner' | 'member' | 'read_only';

export interface Project {
    id: string;
    title: string;
    description: string;
    status: ProjectStatus;
    priority: Priority;
    role: ProjectRole;
    createdAt: string;
}