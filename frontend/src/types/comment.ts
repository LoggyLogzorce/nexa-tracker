export interface CommentUser {
    id: string;
    name: string;
    email: string;
    avatar_url?: string;
}

export interface Comment {
    id: number;
    created_at: string;
    user: CommentUser;
    task_id: number;
    content: string;
}
