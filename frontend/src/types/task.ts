export interface TaskStatus {
    id: number;
    name: string;
    color: string;
    order_index: number;
}

export interface TaskPriority {
    id: number;
    title: string;
    color: string;
}

export interface TaskUser {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
}

export interface Task {
    id: number;
    created_at: string;
    updated_at: string;
    title: string;
    description: string;
    deadline: string | null;
    project_id: string;
    project_title?: string;
    status: TaskStatus;
    priority: TaskPriority;
    assignee: TaskUser | null;
    reporter: TaskUser;
    is_archive: boolean;
}

export interface ProjectMember {
    project_id: string;
    role: string;
    User: {
        user_id: string;
        name: string;
        email: string;
        avatar_url?: string;
    };
}

export interface Attachment {
    id: number;
    created_at: string;
    task_id: number;
    user: TaskUser;
    filename: string;
    file_size: number;
    mime_type: string;
}
