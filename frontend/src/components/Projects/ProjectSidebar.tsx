import type { ProjectMember, Attachment } from '../../types/task';
import TeamMembers from './TeamMembers';
import Attachments from './Attachments';
import styles from './ProjectSidebar.module.css';

interface Props {
    members: ProjectMember[];
    attachments: Attachment[];
    open: boolean;
    projectId: string;
    onAddMember?: () => void;
    onRoleChange?: (member: ProjectMember, role: string) => void;
    onRemoveMember?: (member: ProjectMember) => void;
}

export default function ProjectSidebar({ members, attachments, open, projectId, onAddMember, onRoleChange, onRemoveMember }: Props) {
    return (
        <aside className={`${styles.sidebar} ${open ? styles.open : ''}`}>
            <TeamMembers members={members} onAddMember={onAddMember} onRoleChange={onRoleChange} onRemoveMember={onRemoveMember} />
            <Attachments attachments={attachments} projectId={projectId} />
        </aside>
    );
}
