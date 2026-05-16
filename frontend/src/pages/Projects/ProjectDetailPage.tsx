import { useState, useEffect, useCallback, useMemo } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getProjectById, getProjectTasks, getProjectMembers, getProjectAttachments, updateProjectApi, deleteProjectApi, createTaskApi, updateTaskApi } from '../../api/projects';
import type { Task, ProjectMember, Attachment } from '../../types/task';
import type { Project } from '../../types/project';
import { useNotifications } from '../../contexts/useNotifications';
import ProjectHeader from '../../components/Projects/ProjectHeader';
import ViewSwitcher from '../../components/Projects/ViewSwitcher';
import FilterPanel from '../../components/Projects/FilterPanel';
import KanbanBoard from '../../components/Projects/KanbanBoard';
import ListView from '../../components/Projects/ListView';
import ProjectSidebar from '../../components/Projects/ProjectSidebar';
import NewTaskModal from '../../components/Projects/NewTaskModal';
import EditModal from '../../components/Dashboard/EditModal';
import DeleteModal from '../../components/Dashboard/DeleteModal';
import EditStatusesModal from '../../components/Projects/EditStatusesModal';
import EditPrioritiesModal from '../../components/Projects/EditPrioritiesModal';
import AddMemberModal from '../../components/Projects/AddMemberModal';
import styles from './ProjectDetailPage.module.css';

type View = 'kanban' | 'list';

export default function ProjectDetailPage() {
    const { id } = useParams<{ id: string }>();
    const navigate = useNavigate();
    const { addNotification } = useNotifications();

    const [project, setProject] = useState<Project | null>(null);
    const [tasks, setTasks] = useState<Task[]>([]);
    const [members, setMembers] = useState<ProjectMember[]>([]);
    const [attachments, setAttachments] = useState<Attachment[]>([]);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [view, setView] = useState<View>('kanban');
    const [filteredTasks, setFilteredTasks] = useState<Task[]>([]);
    const [sidebarOpen, setSidebarOpen] = useState(false);
    const [showNewTask, setShowNewTask] = useState(false);
    const [editModal, setEditModal] = useState<{ open: boolean }>({ open: false });
    const [deleteModal, setDeleteModal] = useState<{ open: boolean }>({ open: false });
    const [statusesModal, setStatusesModal] = useState(false);
    const [prioritiesModal, setPrioritiesModal] = useState(false);
    const [showAddMember, setShowAddMember] = useState(false);
    const [filters, setFilters] = useState({ status: '', priority: '', assignee: '' });

    useEffect(() => {
        if (!id) return;

        Promise.all([
            getProjectById(id),
            getProjectTasks(id),
            getProjectMembers(id),
            getProjectAttachments(id),
        ])
            .then(([p, t, m, a]) => {
                setProject(p);
                setTasks(t);
                setFilteredTasks(t);
                setMembers(m);
                setAttachments(a);
                setError(null);
            })
            .catch(err => setError(err instanceof Error ? err.message : 'Ошибка загрузки'))
            .finally(() => setIsLoading(false));
    }, [id]);

    const allMembers = useMemo(() => {
        if (!project) return members;
        const ownerMember: ProjectMember = {
            project_id: project.id,
            role: 'owner',
            User: {
                user_id: project.owner.id,
                name: project.owner.name,
                email: project.owner.email,
            },
        };
        return [ownerMember, ...members];
    }, [project, members]);

    const applyFilters = useCallback((f: { status: string; priority: string; assignee: string }) => {
        const filtered = tasks.filter(task => {
            if (f.status && task.status.name !== f.status) return false;
            if (f.priority && task.priority.title !== f.priority) return false;
            if (f.assignee && (!task.assignee || task.assignee.id !== f.assignee)) return false;
            return true;
        });
        setFilteredTasks(filtered);
    }, [tasks]);

    const handleFilterChange = useCallback((f: { status: string; priority: string; assignee: string }) => {
        setFilters(f);
    }, []);

    const handleApplyFilters = useCallback(() => {
        applyFilters(filters);
    }, [applyFilters, filters]);

    const handleClearFilters = useCallback(() => {
        setFilters({ status: '', priority: '', assignee: '' });
        setFilteredTasks(tasks);
    }, [tasks]);

    const handleRemoveFilter = useCallback((key: 'status' | 'priority' | 'assignee') => {
        const next = { ...filters, [key]: '' };
        setFilters(next);
        applyFilters(next);
    }, [filters, applyFilters]);

    const handleCreateTask = useCallback((data: { title: string; description: string; status: string; priority: string; deadline: string; assignee: string }) => {
        if (!project || !project.statuses || !project.priorities) return;

        const statusObj = project.statuses.find(s => s.name === data.status) || project.statuses[0];
        const priorityObj = project.priorities.find(p => p.title === data.priority) || project.priorities[0];
        const assigneeUser = data.assignee ? allMembers.find(m => m.User.user_id === data.assignee) : null;

        createTaskApi(project.id, {
            title: data.title,
            description: data.description,
            deadline: data.deadline ? `${data.deadline}T00:00:00Z` : null,
            status_id: statusObj.id,
            priority_id: priorityObj.id,
            assignee_id: assigneeUser?.User.user_id ?? null,
        })
            .then(task => {
                setTasks(prev => [...prev, task]);
                setFilteredTasks(prev => [...prev, task]);
                setShowNewTask(false);
                addNotification('success', 'Задача успешно создана');
            })
            .catch(() => addNotification('error', 'Ошибка при создании задачи'));
    }, [project, allMembers, addNotification]);

    const handleSaveEdit = useCallback((updated: Project) => {
        if (!project) return;
        updateProjectApi(updated.id, updated)
            .then(p => { setProject(p); setEditModal({ open: false }); addNotification('success', 'Проект успешно обновлён'); })
            .catch(() => addNotification('error', 'Ошибка при обновлении проекта'));
    }, [project, addNotification]);

    const handleConfirmDelete = useCallback(() => {
        if (!project) return;
        deleteProjectApi(project.id)
            .then(() => { addNotification('success', 'Проект успешно удалён'); navigate('/projects'); })
            .catch(() => addNotification('error', 'Ошибка при удалении проекта'));
    }, [project, addNotification, navigate]);

    const handleMemberAdded = useCallback(() => {
        if (!id) return;
        getProjectMembers(id).then(m => setMembers(m));
        setShowAddMember(false);
    }, [id]);

    const handleTaskMove = useCallback((taskId: number, newStatusName: string) => {
        if (!project) return;
        const newStatus = project.statuses.find(s => s.name === newStatusName);
        if (!newStatus) return;

        setTasks(prev => prev.map(t => t.id === taskId ? { ...t, status: newStatus } : t));
        setFilteredTasks(prev => prev.map(t => t.id === taskId ? { ...t, status: newStatus } : t));
        updateTaskApi(project.id, taskId, { status_id: newStatus.id })
            .catch(() => addNotification('error', 'Ошибка при перемещении задачи'));
    }, [project, addNotification]);

    const handleRefetchProject = useCallback(() => {
        if (!id) return;
        getProjectById(id).then(p => setProject(p));
    }, [id]);

    if (isLoading) {
        return <p className={styles.loading}>Загрузка проекта...</p>;
    }

    if (error || !project) {
        return (
            <div className={styles.errorBlock}>
                <p className={styles.error}>{error || 'Проект не найден'}</p>
            </div>
        );
    }

    return (
        <div className={`${styles.page} ${sidebarOpen ? styles.hasSidebar : ''}`}>
            <div className={`${styles.center} ${sidebarOpen ? styles.centerWithSidebar : ''}`}>
                <ProjectHeader project={project} onEdit={() => setEditModal({ open: true })} onDelete={() => setDeleteModal({ open: true })} onEditStatuses={() => setStatusesModal(true)} onEditPriorities={() => setPrioritiesModal(true)} />

                <div className={styles.toolbar}>
                    <div className={styles.toolbarLeft}>
                        <ViewSwitcher active={view} onChange={setView} />
                        <div className={styles.divider} />
                        <FilterPanel
                            statuses={project.statuses}
                            priorities={project.priorities}
                            members={allMembers}
                            filters={filters}
                            onFilterChange={handleFilterChange}
                            onApply={handleApplyFilters}
                            onClear={handleClearFilters}
                        />
                        {filters.status && (
                            <span className={styles.chip}>
                                {filters.status}
                                <button className={styles.chipRemove} onClick={() => handleRemoveFilter('status')}>
                                    <svg className={styles.chipRemoveIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
                                </button>
                            </span>
                        )}
                        {filters.priority && (
                            <span className={styles.chip}>
                                {filters.priority}
                                <button className={styles.chipRemove} onClick={() => handleRemoveFilter('priority')}>
                                    <svg className={styles.chipRemoveIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
                                </button>
                            </span>
                        )}
                        {filters.assignee && (
                            <span className={styles.chip}>
                                {allMembers.find(m => m.User.user_id === filters.assignee)?.User.name || filters.assignee}
                                <button className={styles.chipRemove} onClick={() => handleRemoveFilter('assignee')}>
                                    <svg className={styles.chipRemoveIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12"/></svg>
                                </button>
                            </span>
                        )}
                    </div>
                    <button className={styles.newTaskBtn} onClick={() => setShowNewTask(true)}>
                        <svg className={styles.newTaskIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                        Новая задача
                    </button>
                </div>

                {view === 'kanban' ? (
                    <KanbanBoard statuses={project.statuses} tasks={filteredTasks} onTaskMove={handleTaskMove} />
                ) : (
                    <ListView tasks={filteredTasks} />
                )}
            </div>

            <button className={styles.floatingToggle} onClick={() => setSidebarOpen(!sidebarOpen)}>
                <svg className={styles.floatingToggleIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d={sidebarOpen ? 'M15 19l-7-7 7-7' : 'M9 5l7 7-7 7'} /></svg>
            </button>

            <ProjectSidebar members={allMembers} attachments={attachments} open={sidebarOpen} onAddMember={() => setShowAddMember(true)} />

            {showNewTask && (
                <NewTaskModal
                    statuses={project.statuses}
                    priorities={project.priorities}
                    members={allMembers}
                    onClose={() => setShowNewTask(false)}
                    onSave={handleCreateTask}
                />
            )}

            {editModal.open && project && (
                <EditModal project={project} onClose={() => setEditModal({ open: false })} onSave={handleSaveEdit} />
            )}

            {deleteModal.open && project && (
                <DeleteModal project={project} onClose={() => setDeleteModal({ open: false })} onConfirm={() => handleConfirmDelete()} />
            )}

            {statusesModal && project && (
                <EditStatusesModal projectId={project.id} statuses={project.statuses} onClose={() => setStatusesModal(false)} onSave={() => { setStatusesModal(false); handleRefetchProject(); }} />
            )}

            {prioritiesModal && project && (
                <EditPrioritiesModal projectId={project.id} priorities={project.priorities} onClose={() => setPrioritiesModal(false)} onSave={() => { setPrioritiesModal(false); handleRefetchProject(); }} />
            )}

            {showAddMember && project && (
                <AddMemberModal projectId={project.id} onClose={() => setShowAddMember(false)} onSave={handleMemberAdded} />
            )}
        </div>
    );
}
