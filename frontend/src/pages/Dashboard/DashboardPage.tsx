import { useState } from 'react';
import type {Project} from '../../types/project';
import { useProjects } from '../../hooks/useProjects';
import ProjectCard from '../../components/Dashboard/ProjectCard';
import EditModal from '../../components/Dashboard/EditModal';
import DeleteModal from '../../components/Dashboard/DeleteModal';
import CreateModal from '../../components/Projects/CreateModal';
import styles from './DashboardPage.module.css';

export default function DashboardPage() {
    const { myProjects, createProject, updateProject, deleteProject } = useProjects();
    const [showCreate, setShowCreate] = useState(false);
    const [editModal, setEditModal] = useState<{ open: boolean; project?: Project }>({ open: false });
    const [deleteModal, setDeleteModal] = useState<{ open: boolean; project?: Project }>({ open: false });

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

    return (
        <div className={styles.page}>
            <div className={styles.header}>
                <div>
                    <h1 className={styles.title}>Мои Проекты</h1>
                    <p className={styles.subtitle}>Управление активными задачами команды</p>
                </div>
                <button className={styles.newBtn} onClick={() => setShowCreate(true)}>
                    <svg className={styles.icon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                    Новый проект
                </button>
            </div>

            <div className={styles.grid}>
                {myProjects.map(project => (
                    <ProjectCard key={project.id} project={project} onEdit={handleEdit} onDelete={handleDelete} />
                ))}
                <button className={styles.addCard} onClick={() => setShowCreate(true)}>
                    <div className={styles.addCardContent}>
                        <svg className={styles.addIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                        <span>Создать новый проект</span>
                    </div>
                </button>
            </div>

            {showCreate && (
                <CreateModal onClose={() => setShowCreate(false)} onSave={createProject} />
            )}
            {editModal.open && editModal.project && (
                <EditModal project={editModal.project} onClose={() => setEditModal({ open: false })} onSave={handleSaveEdit} />
            )}
            {deleteModal.open && deleteModal.project && (
                <DeleteModal project={deleteModal.project} onClose={() => setDeleteModal({ open: false })} onConfirm={handleConfirmDelete} />
            )}
        </div>
    );
}