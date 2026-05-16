import { useDroppable } from '@dnd-kit/core';
import type { TaskStatus, Task } from '../../types/task';
import KanbanCard from './KanbanCard';
import styles from './KanbanColumn.module.css';

interface Props { status: TaskStatus; tasks: Task[]; }

export default function KanbanColumn({ status, tasks }: Props) {
    const { setNodeRef, isOver } = useDroppable({ id: status.name });

    return (
        <div className={styles.column}>
            <div className={styles.header}>
                <h3 className={styles.title}>
                    <span className={styles.dot} style={{ backgroundColor: status.color }} />
                    {status.name}
                    <span className={styles.count}>{tasks.length}</span>
                </h3>
            </div>
            <div ref={setNodeRef} className={`${styles.list} ${isOver ? styles.dropOver : ''}`}>
                {tasks.map(task => (
                    <KanbanCard key={task.id} task={task} />
                ))}
            </div>
        </div>
    );
}
