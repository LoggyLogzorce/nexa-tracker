export type ProjectStatus = 'В работе' | 'Планирование' | 'Завершен';
export type Priority = 'Низкий' | 'Средний' | 'Высокий';
export type ProjectRole = 'owner' | 'member' | 'read_only';

export interface ProjectOwner {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
}

export interface CustomStatus {
    id: number;
    name: string;
    color: string;
    order_index: number;
}

export interface CustomPriority {
    id: number;
    title: string;
    color: string;
}

export interface Project {
    id: string;
    title: string;
    description: string;
    status: ProjectStatus;
    priority: Priority;
    role: ProjectRole;
    owner: ProjectOwner;
    createdAt: string;
    statuses: CustomStatus[];
    priorities: CustomPriority[];
}