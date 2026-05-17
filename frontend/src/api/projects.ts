import client, { extractData } from './client';
import type { Project, ProjectRole, ProjectOwner, ProjectStatus, Priority, CustomStatus, CustomPriority } from '../types/project';
import type { Task, ProjectMember, Attachment } from '../types/task';
import type { Comment } from '../types/comment';
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

export const getProjectTasks = async (projectId: string, archived?: boolean): Promise<Task[]> => {
    const response = await client.get<ApiResponse<Task[]>>(`/projects/${projectId}/tasks`, {
        params: archived ? { archived: 'true' } : undefined,
    });
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
    is_archive?: boolean;
}

export const createTaskApi = async (projectId: string, data: CreateTaskPayload): Promise<Task> => {
    const response = await client.post<ApiResponse<Task>>(`/projects/${projectId}/tasks`, data);
    return extractData(response.data);
};

export const updateTaskApi = async (projectId: string, taskId: number, data: Partial<CreateTaskPayload>, archived?: boolean): Promise<Task> => {
    const response = await client.put<ApiResponse<Task>>(`/projects/${projectId}/tasks/${taskId}`, data, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const getTaskById = async (projectId: string, taskId: number, archived?: boolean): Promise<Task> => {
    const response = await client.get<ApiResponse<Task>>(`/projects/${projectId}/tasks/${taskId}`, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const getTaskComments = async (projectId: string, taskId: number, archived?: boolean): Promise<Comment[]> => {
    const response = await client.get<ApiResponse<Comment[]>>(`/projects/${projectId}/tasks/${taskId}/comments`, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const createTaskComment = async (projectId: string, taskId: number, content: string, archived?: boolean): Promise<Comment> => {
    const response = await client.post<ApiResponse<Comment>>(`/projects/${projectId}/tasks/${taskId}/comments`, { content }, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const getTaskAttachments = async (projectId: string, taskId: number, archived?: boolean): Promise<Attachment[]> => {
    const response = await client.get<ApiResponse<Attachment[]>>(`/projects/${projectId}/tasks/${taskId}/attachments`, {
        params: archived ? { archived: 'true' } : undefined,
    });
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

export const updateMemberRoleApi = async (projectId: string, userId: string, role: string): Promise<void> => {
    await client.put(`/projects/${projectId}/participants/${userId}`, { role });
};

export const removeMemberApi = async (projectId: string, userId: string): Promise<void> => {
    await client.delete(`/projects/${projectId}/participants/${userId}`);
};

// Task deletion
export const deleteTaskApi = async (projectId: string, taskId: number, archived?: boolean): Promise<void> => {
    await client.delete(`/projects/${projectId}/tasks/${taskId}`, {
        params: archived ? { archived: 'true' } : undefined,
    });
};

// Task comments update/delete
export const updateTaskCommentApi = async (projectId: string, taskId: number, commentId: number, content: string, archived?: boolean): Promise<Comment> => {
    const response = await client.put<ApiResponse<Comment>>(`/projects/${projectId}/tasks/${taskId}/comments/${commentId}`, { content }, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const deleteTaskCommentApi = async (projectId: string, taskId: number, commentId: number, archived?: boolean): Promise<void> => {
    await client.delete(`/projects/${projectId}/tasks/${taskId}/comments/${commentId}`, {
        params: archived ? { archived: 'true' } : undefined,
    });
};

// Task attachments upload/delete/download
export const uploadTaskAttachmentApi = async (projectId: string, taskId: number, file: File, archived?: boolean): Promise<Attachment> => {
    const formData = new FormData();
    formData.append('file', file);
    const response = await client.post<ApiResponse<Attachment>>(`/projects/${projectId}/tasks/${taskId}/attachments`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

export const deleteTaskAttachmentApi = async (projectId: string, taskId: number, attachmentId: number, archived?: boolean): Promise<void> => {
    await client.delete(`/projects/${projectId}/tasks/${taskId}/attachments/${attachmentId}`, {
        params: archived ? { archived: 'true' } : undefined,
    });
};

export const downloadAttachmentApi = async (projectId: string, taskId: number, attachmentId: number, filename: string, archived?: boolean): Promise<void> => {
    const response = await client.get(`/projects/${projectId}/tasks/${taskId}/attachments/${attachmentId}`, {
        responseType: 'blob',
        params: archived ? { archived: 'true' } : undefined,
    });
    const url = window.URL.createObjectURL(new Blob([response.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
    window.URL.revokeObjectURL(url);
};

interface HistoryUser {
    id: string; name: string; email: string; avatar_url?: string;
}

export interface TaskHistoryEntry {
    id: number;
    created_at: string;
    user: HistoryUser;
    task_id: number;
    old: Record<string, unknown>;
    new: Record<string, unknown>;
    changes: Array<{ field: string; old_value: string; new_value: string }>;
}

export const getTaskHistoryApi = async (projectId: string, taskId: number, archived?: boolean): Promise<TaskHistoryEntry[]> => {
    const response = await client.get<ApiResponse<TaskHistoryEntry[]>>(`/projects/${projectId}/tasks/${taskId}/history`, {
        params: archived ? { archived: 'true' } : undefined,
    });
    return extractData(response.data);
};

// Tasks by user (assigned/reported)
export const getTasksByUserApi = async (type: 'assigned' | 'reported'): Promise<Task[]> => {
    const response = await client.get<ApiResponse<Task[]>>(`/tasks/me?type=${type}`);
    return extractData(response.data);
};

// Search
export const searchTasksApi = async (q: string): Promise<Task[]> => {
    const response = await client.get<ApiResponse<Task[]>>(`/tasks/search?q=${encodeURIComponent(q)}`);
    return extractData(response.data);
};

export const searchProjectsApi = async (q: string): Promise<Project[]> => {
    const response = await client.get<ApiResponse<Project[]>>(`/projects/search?q=${encodeURIComponent(q)}`);
    return extractData(response.data);
};

// Standalone statuses and priorities (if needed separately from project)
export const getProjectStatusesApi = async (projectId: string): Promise<CustomStatus[]> => {
    const response = await client.get<ApiResponse<CustomStatus[]>>(`/projects/${projectId}/statuses`);
    return extractData(response.data);
};

export const getProjectPrioritiesApi = async (projectId: string): Promise<CustomPriority[]> => {
    const response = await client.get<ApiResponse<CustomPriority[]>>(`/projects/${projectId}/priorities`);
    return extractData(response.data);
};
