import type { ProjectMember, Attachment } from '../../types/task';
import TeamMembers from './TeamMembers';
import Attachments from './Attachments';
import styles from './ProjectSidebar.module.css';

interface Props {
    members: ProjectMember[];
    attachments: Attachment[];
    open: boolean;
    onAddMember?: () => void;
}

export default function ProjectSidebar({ members, attachments, open, onAddMember }: Props) {
    return (
        <aside className={`${styles.sidebar} ${open ? styles.open : ''}`}>
            <TeamMembers members={members} onAddMember={onAddMember} />
            <Attachments attachments={attachments} />
        </aside>
    );
}
