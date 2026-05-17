import { useState } from 'react';
import type { ProjectMember } from '../../types/task';
import Avatar from '../UI/Avatar';
import styles from './TeamMembers.module.css';

interface Props { members: ProjectMember[]; onAddMember?: () => void; onRoleChange?: (member: ProjectMember, role: string) => void; onRemoveMember?: (member: ProjectMember) => void; }

export default function TeamMembers({ members, onAddMember, onRoleChange, onRemoveMember }: Props) {
    const [openMenu, setOpenMenu] = useState<string | null>(null);
    const [roleSubmenu, setRoleSubmenu] = useState(false);

    return (
        <div className={styles.wrap}>
            <div className={styles.header}>
                <h3 className={styles.title}>Команда проекта</h3>
                <button className={styles.addBtn} onClick={onAddMember}>
                    <svg className={styles.addIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4"/></svg>
                </button>
            </div>
            <div className={styles.list}>
                {members.filter(m => m.User.name !== 'Deleted User').map(member => (
                    <div key={member.User.user_id} className={styles.member}>
                        <div className={styles.memberLeft}>
                            <Avatar name={member.User.name} avatarUrl={member.User.avatar_url} size={32} />
                            <div className={styles.info}>
                                <p className={styles.name}>{member.User.name}</p>
                                <p className={styles.role}>{member.role === 'owner' ? 'Владелец' : member.role === 'read_only' ? 'Только чтение' : 'Участник'}</p>
                            </div>
                        </div>
                        {member.role !== 'owner' && (
                            <div className={styles.menuWrap}>
                                <button className={styles.menuBtn} onClick={() => { setOpenMenu(openMenu === member.User.user_id ? null : member.User.user_id); setRoleSubmenu(false); }}>
                                    <svg className={styles.menuIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z"/></svg>
                                </button>
                                {openMenu === member.User.user_id && !roleSubmenu && (
                                    <div className={styles.dropdown}>
                                        <button className={styles.dropdownItem} onClick={() => setRoleSubmenu(true)}>
                                            <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/></svg>
                                            Роль
                                        </button>
                                        <button className={`${styles.dropdownItem} ${styles.dropdownDelete}`} onClick={() => { onRemoveMember?.(member); setOpenMenu(null); setRoleSubmenu(false); }}>
                                            <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"/></svg>
                                            Удалить
                                        </button>
                                    </div>
                                )}
                                {openMenu === member.User.user_id && roleSubmenu && (
                                    <div className={styles.dropdown}>
                                        <button className={styles.dropdownItem} onClick={() => setRoleSubmenu(false)}>
                                            <svg className={styles.dropdownIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7"/></svg>
                                            Назад
                                        </button>
                                        <div className={styles.submenuDivider} />
                                        <button className={`${styles.dropdownItem} ${member.role === 'member' ? styles.activeRole : ''}`} onClick={() => { onRoleChange?.(member, 'member'); setOpenMenu(null); setRoleSubmenu(false); }}>
                                            Участник
                                        </button>
                                        <button className={`${styles.dropdownItem} ${member.role === 'read_only' ? styles.activeRole : ''}`} onClick={() => { onRoleChange?.(member, 'read_only'); setOpenMenu(null); setRoleSubmenu(false); }}>
                                            Только чтение
                                        </button>
                                    </div>
                                )}
                            </div>
                        )}
                    </div>
                ))}
            </div>
        </div>
    );
}
