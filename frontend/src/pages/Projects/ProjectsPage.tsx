import { useState } from 'react';
import type { Project } from '../../types/project';
import { useProjects } from '../../hooks/useProjects';
import { getProjectById } from '../../api/projects';
import ProjectCard from '../../components/Dashboard/ProjectCard';
import EditModal from '../../components/Dashboard/EditModal';
import DeleteModal from '../../components/Dashboard/DeleteModal';
import CreateModal from '../../components/Projects/CreateModal';
import EditStatusesModal from '../../components/Projects/EditStatusesModal';
import EditPrioritiesModal from '../../components/Projects/EditPrioritiesModal';
import styles from './ProjectsPage.module.css';

export default function ProjectsPage() {
    const { projects, isLoading, createProject, updateProject, deleteProject } = useProjects('all');
    const [showCreate, setShowCreate] = useState(false);
    const [editModal, setEditModal] = useState<{ open: boolean; project?: Project }>({ open: false });
    const [deleteModal, setDeleteModal] = useState<{ open: boolean; project?: Project }>({ open: false });
    const [statusEdit, setStatusEdit] = useState<{ open: boolean; project?: Project }>({ open: false });
    const [priorityEdit, setPriorityEdit] = useState<{ open: boolean; project?: Project }>({ open: false });

    const handleEdit = (project: Project) => { setEditModal({ open: true, project }); };
    const handleDelete = (project: Project) => { setDeleteModal({ open: true, project }); };

    const handleSaveEdit = (updated: Project) => {
        updateProject(updated);
        setEditModal({ open: false });
    };

    const handleConfirmDelete = (id: string) => {
        deleteProject(id);
        setDeleteModal({ open: false });
    };

    const handleSaveStatuses = (project: Project) => {
        getProjectById(project.id).then(p => { updateProject(p); setStatusEdit({ open: false }); });
    };

    const handleSavePriorities = (project: Project) => {
        getProjectById(project.id).then(p => { updateProject(p); setPriorityEdit({ open: false }); });
    };

    return (
        <div className={styles.page}>
            <div className={styles.header}>
                <div>
                    <h1 className={styles.title}>Проекты</h1>
                    <p className={styles.subtitle}>Все проекты в которых вы участвуете</p>
                </div>
                <button className={styles.newBtn} onClick={() => setShowCreate(true)}>
                    <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                    Новый проект
                </button>
            </div>

            {isLoading ? (
                <p className={styles.loading}>Загрузка проектов...</p>
            ) : (
            <div className={styles.grid}>
                {projects.map(project => (
                    <ProjectCard key={project.id} project={project} onEdit={handleEdit} onDelete={handleDelete} onEditStatuses={p => setStatusEdit({ open: true, project: p })} onEditPriorities={p => setPriorityEdit({ open: true, project: p })} />
                ))}
                <button className={styles.addCard} onClick={() => setShowCreate(true)}>
                    <div className={styles.addCardContent}>
                        <svg className={styles.addIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                        <span>Создать новый проект</span>
                    </div>
                </button>
            </div>
            )}

            {showCreate && (
                <CreateModal onClose={() => setShowCreate(false)} onSave={createProject} />
            )}
            {editModal.open && editModal.project && (
                <EditModal project={editModal.project} onClose={() => setEditModal({ open: false })} onSave={handleSaveEdit} />
            )}
            {deleteModal.open && deleteModal.project && (
                <DeleteModal project={deleteModal.project} onClose={() => setDeleteModal({ open: false })} onConfirm={handleConfirmDelete} />
            )}

            {statusEdit.open && statusEdit.project && (
                <EditStatusesModal projectId={statusEdit.project.id} statuses={statusEdit.project.statuses} onClose={() => setStatusEdit({ open: false })} onSave={() => handleSaveStatuses(statusEdit.project!)} />
            )}

            {priorityEdit.open && priorityEdit.project && (
                <EditPrioritiesModal projectId={priorityEdit.project.id} priorities={priorityEdit.project.priorities} onClose={() => setPriorityEdit({ open: false })} onSave={() => handleSavePriorities(priorityEdit.project!)} />
            )}
        </div>
    );
}
