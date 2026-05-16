import { DndContext, DragOverlay, PointerSensor, useSensor, useSensors, type DragEndEvent, type DragStartEvent } from '@dnd-kit/core';
import { useState } from 'react';
import type { TaskStatus, Task } from '../../types/task';
import KanbanColumn from './KanbanColumn';
import KanbanCard from './KanbanCard';
import styles from './KanbanBoard.module.css';

interface Props { statuses: TaskStatus[]; tasks: Task[]; onTaskMove?: (taskId: number, newStatusName: string) => void; }

export default function KanbanBoard({ statuses, tasks, onTaskMove }: Props) {
    const [activeTask, setActiveTask] = useState<Task | null>(null);
    const sorted = [...statuses].sort((a, b) => a.order_index - b.order_index);

    const sensors = useSensors(useSensor(PointerSensor, { activationConstraint: { distance: 8 } }));

    const handleDragStart = (event: DragStartEvent) => {
        setActiveTask(tasks.find(t => String(t.id) === event.active.id) ?? null);
    };

    const handleDragEnd = (event: DragEndEvent) => {
        setActiveTask(null);
        const { active, over } = event;
        if (!over || active.id === over.id) return;
        const taskId = Number(active.id);
        const newStatusName = String(over.id);
        onTaskMove?.(taskId, newStatusName);
    };

    return (
        <DndContext sensors={sensors} onDragStart={handleDragStart} onDragEnd={handleDragEnd}>
            <div className={styles.board}>
                {sorted.map(status => (
                    <KanbanColumn
                        key={status.id}
                        status={status}
                        tasks={tasks.filter(t => t.status.name === status.name)}
                    />
                ))}
            </div>
            <DragOverlay>
                {activeTask ? <div className={styles.overlay}><KanbanCard task={activeTask} /></div> : null}
            </DragOverlay>
        </DndContext>
    );
}
