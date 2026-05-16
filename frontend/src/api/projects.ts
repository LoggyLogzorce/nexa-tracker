import client, { extractData } from './client';
import type { Project, ProjectRole, ProjectOwner, ProjectStatus, Priority, CustomStatus, CustomPriority } from '../types/project';
import type { Task, ProjectMember, Attachment } from '../types/task';
import type { ApiResponse } from '../types/auth';

interface RawProject {
    id: string;
    title: string;
    description: string | null;
    created_at: string;
    owner: ProjectOwner;
    status: string;
    priority: string;
    user_role: ProjectRole;
    statuses: CustomStatus[];
    priorities: CustomPriority[];
}

const statusToApi: Record<ProjectStatus, string> = {
    'Планирование': 'plan',
    'В работе': 'in_progress',
    'Завершен': 'done',
};

const priorityToApi: Record<Priority, string> = {
    'Низкий': 'low',
    'Средний': 'medium',
    'Высокий': 'high',
};

const statusFromApi: Record<string, ProjectStatus> = {
    plan: 'Планирование',
    in_progress: 'В работе',
    done: 'Завершен',
};

const priorityFromApi: Record<string, Priority> = {
    low: 'Низкий',
    medium: 'Средний',
    high: 'Высокий',
};

function mapProject(raw: RawProject): Project {
    return {
        id: raw.id,
        title: raw.title,
        description: raw.description ?? '',
        status: statusFromApi[raw.status] ?? 'Планирование',
        priority: priorityFromApi[raw.priority] ?? 'Средний',
        role: raw.user_role,
        owner: raw.owner,
        createdAt: raw.created_at,
        statuses: raw.statuses ?? [],
        priorities: raw.priorities ?? [],
    };
}

export const getProjects = async (): Promise<Project[]> => {
    const response = await client.get<ApiResponse<RawProject[]>>('/projects');
    return extractData(response.data).map(mapProject);
};

export const getOwnedProjects = async (): Promise<Project[]> => {
    const response = await client.get<ApiResponse<RawProject[]>>('/projects/owned');
    return extractData(response.data).map(mapProject);
};

export const getProjectById = async (id: string): Promise<Project> => {
    const response = await client.get<ApiResponse<RawProject>>(`/projects/${id}`);
    return mapProject(extractData(response.data));
};

export interface CreateProjectPayload {
    title: string;
    description: string;
    status: string;
    priority: string;
}

export const createProjectApi = async (data: { title: string; description: string; status: ProjectStatus; priority: Priority }): Promise<Project> => {
    const payload: CreateProjectPayload = {
        title: data.title,
        description: data.description,
        status: statusToApi[data.status],
        priority: priorityToApi[data.priority],
    };
    const response = await client.post<ApiResponse<RawProject>>('/projects', payload);
    return mapProject(extractData(response.data));
};

export const updateProjectApi = async (id: string, data: { title: string; description: string; status: ProjectStatus; priority: Priority }): Promise<Project> => {
    const payload: CreateProjectPayload = {
        title: data.title,
        description: data.description,
        status: statusToApi[data.status],
        priority: priorityToApi[data.priority],
    };
    const response = await client.put<ApiResponse<RawProject>>(`/projects/${id}`, payload);
    return mapProject(extractData(response.data));
};

export const deleteProjectApi = async (id: string): Promise<void> => {
    await client.delete(`/projects/${id}`);
};

export const getProjectTasks = async (projectId: string): Promise<Task[]> => {
    const response = await client.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks`);
    return extractData(response.data);
};

export const getProjectMembers = async (projectId: string): Promise<ProjectMember[]> => {
    const response = await client.get<ApiResponse<ProjectMember[]>>(`/projects/${projectId}/participants`);
    return extractData(response.data);
};

export const getProjectAttachments = async (projectId: string): Promise<Attachment[]> => {
    const response = await client.get<ApiResponse<Attachment[]>>(`/projects/${projectId}/attachments`);
    return extractData(response.data);
};

// Statuses
export const createStatusApi = async (projectId: string, data: { name: string; color: string; order_index: number }): Promise<CustomStatus> => {
    const response = await client.post<ApiResponse<CustomStatus>>(`/projects/${projectId}/statuses`, data);
    return extractData(response.data);
};

export const updateStatusApi = async (projectId: string, statusId: number, data: { name: string; color: string; order_index: number }): Promise<CustomStatus> => {
    const response = await client.put<ApiResponse<CustomStatus>>(`/projects/${projectId}/statuses/${statusId}`, data);
    return extractData(response.data);
};

export const deleteStatusApi = async (projectId: string, statusId: number): Promise<void> => {
    await client.delete(`/projects/${projectId}/statuses/${statusId}`);
};

// Priorities
export const createPriorityApi = async (projectId: string, data: { title: string; color: string }): Promise<CustomPriority> => {
    const response = await client.post<ApiResponse<CustomPriority>>(`/projects/${projectId}/priorities`, data);
    return extractData(response.data);
};

export const updatePriorityApi = async (projectId: string, priorityId: number, data: { title: string; color: string }): Promise<CustomPriority> => {
    const response = await client.put<ApiResponse<CustomPriority>>(`/projects/${projectId}/priorities/${priorityId}`, data);
    return extractData(response.data);
};

export const deletePriorityApi = async (projectId: string, priorityId: number): Promise<void> => {
    await client.delete(`/projects/${projectId}/priorities/${priorityId}`);
};

// Tasks
export interface CreateTaskPayload {
    title: string;
    description?: string;
    deadline?: string | null;
    status_id?: number;
    priority_id?: number;
    assignee_id?: string | null;
}

export const createTaskApi = async (projectId: string, data: CreateTaskPayload): Promise<Task> => {
    const response = await client.post<ApiResponse<Task>>(`/projects/${projectId}/tasks`, data);
    return extractData(response.data);
};

export const updateTaskApi = async (projectId: string, taskId: number, data: Partial<CreateTaskPayload>): Promise<Task> => {
    const response = await client.put<ApiResponse<Task>>(`/projects/${projectId}/tasks/${taskId}`, data);
    return extractData(response.data);
};

// Users
export interface UserSearchResult {
    id: string;
    name: string;
    email: string;
}

export const searchUsersApi = async (email: string): Promise<UserSearchResult[]> => {
    const response = await client.get<ApiResponse<UserSearchResult[]>>(`/users/search?q=${encodeURIComponent(email)}`);
    return extractData(response.data);
};

export const addMemberApi = async (projectId: string, userId: string, role: string = 'member'): Promise<void> => {
    await client.post(`/projects/${projectId}/participants`, { user_id: userId, role });
};
