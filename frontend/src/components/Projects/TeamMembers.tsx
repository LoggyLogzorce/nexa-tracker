import type { ProjectMember } from '../../types/task';
import styles from './TeamMembers.module.css';

const avatarImages = [
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1281ab55e-b3ab-4e9e-be91-85c89891a3c3.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/15745d149-3792-4261-a5c2-f7913b692184.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1818d7289-0020-4e2c-bba4-19c272b7c3ba.png',
    'https://image.qwenlm.ai/public_source/87cc46bb-251e-456e-a7c0-8d3b41f35211/1df69439c-4074-406c-b81d-7d3faea2d322.png'
];

interface Props { members: ProjectMember[]; onAddMember?: () => void; }

export default function TeamMembers({ members, onAddMember }: Props) {
    return (
        <div className={styles.wrap}>
            <div className={styles.header}>
                <h3 className={styles.title}>Команда проекта</h3>
                <button className={styles.addBtn} onClick={onAddMember}>
                    <svg className={styles.addIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                </button>
            </div>
            <div className={styles.list}>
                {members.filter(m => m.User.name !== 'Deleted User').map((member, i) => (
                    <div key={member.User.user_id} className={styles.member}>
                        <div className={styles.memberLeft}>
                            <div className={styles.avatar}>
                                <img src={avatarImages[i % avatarImages.length]} alt={member.User.name} />
                            </div>
                            <div className={styles.info}>
                                <p className={styles.name}>{member.User.name}</p>
                                <p className={styles.role}>{member.role === 'owner' ? 'Владелец' : member.role === 'read_only' ? 'Только чтение' : 'Участник'}</p>
                            </div>
                        </div>
                        <span className={styles.online} />
                    </div>
                ))}
            </div>
        </div>
    );
}
