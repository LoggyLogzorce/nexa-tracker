-- Nexa Task Tracker Database Schema

create table users (
    id uuid PRIMARY KEY UNIQUE,
    email varchar(255) NOT NULL UNIQUE,
    password_hash varchar(255) NOT NULL,
    name varchar(50) NOT NULL,
    role varchar(20) NOT NULL DEFAULT 'user',
    secret_2fa varchar(255) NULL,
    created_at timestamptz not null default now()
);

create table refresh_tokens (
    id serial PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) on delete cascade,
    token_hash varchar(255) NOT NULL UNIQUE,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL default now(),
    revoked_at timestamptz NULL,
    user_agent varchar(255) NULL,
    ip_address inet NULL
);

create table projects (
    id serial PRIMARY KEY,
    title varchar(50) NOT NULL,
    description varchar(255) NULL,
    owner_id uuid NOT NULL REFERENCES users(id) on delete cascade,
    created_at timestamptz default now()
);

create table project_participants (
    project_id int NOT NULL REFERENCES projects(id) on delete cascade,
    user_id uuid NOT NULL REFERENCES users(id) on delete cascade,
    role varchar(10) CHECK ( role in ('owner', 'member', 'read_only')),
    PRIMARY KEY (project_id, user_id)
);

create table statuses (
    id serial PRIMARY KEY,
    project_id int REFERENCES projects(id) on delete cascade,
    name varchar(50) NOT NULL,
    color varchar(7) default '#cccccc',
    order_index int default 0,
    unique (project_id, name)
);

create table priorities (
    id serial PRIMARY KEY,
    project_id int REFERENCES projects(id) on delete cascade,
    title varchar(50) NOT NULL,
    color varchar(7) default '#cccccc',
    unique (project_id, title)
);

create table tasks (
    id serial PRIMARY KEY,
    project_id int REFERENCES projects(id) on delete cascade,
    title varchar(100) NOT NULL,
    description text NULL,
    status_id int NULL REFERENCES statuses(id) on delete set null,
    deadline timestamptz NULL,
    priority_id int NULL REFERENCES priorities(id) on delete set null,
    assignee_id uuid NULL REFERENCES users(id) on delete set null,
    reporter_id uuid NULL REFERENCES users(id) on delete set null,
    created_at timestamptz default now(),
    updated_at timestamptz default now()
);

create table update_history (
    id serial PRIMARY KEY,
    user_id uuid REFERENCES users(id) on delete set null,
    task_id int REFERENCES tasks(id) on delete cascade,
    old jsonb,
    new jsonb,
    created_at timestamptz default now()
);

create table comments (
    id serial PRIMARY KEY,
    user_id uuid REFERENCES users(id) on delete set null,
    task_id int REFERENCES tasks(id) on delete cascade,
    content text NOT NULL,
    created_at timestamptz default now()
);

create table attachments (
    id serial PRIMARY KEY,
    task_id int REFERENCES tasks(id) on delete cascade,
    user_id uuid REFERENCES users(id) on delete set null,
    filename varchar(255) NOT NULL,
    file_path varchar(500) NOT NULL,
    file_size bigint NOT NULL,
    mime_type varchar(100),
    created_at timestamptz default now()
);

-- Indexes for performance
create index idx_refresh_tokens_user_id on refresh_tokens(user_id);
create index idx_refresh_tokens_expires_at on refresh_tokens(expires_at);
create index idx_tasks_project_id on tasks(project_id);
create index idx_tasks_assignee_id on tasks(assignee_id);
create index idx_tasks_status_id on tasks(status_id);
create index idx_tasks_deadline on tasks(deadline);
create index idx_comments_task_id on comments(task_id);
create index idx_update_history_task_id on update_history(task_id);
create index idx_project_participants_user_id on project_participants(user_id);
create index idx_attachments_task_id on attachments(task_id);
