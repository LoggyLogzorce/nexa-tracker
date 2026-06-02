import { useNavigate } from 'react-router-dom';
import styles from './NewVersionModal.module.css';

interface Props {
    version: string;
    onClose: () => void;
}

export default function NewVersionModal({ version, onClose }: Props) {
    const navigate = useNavigate();

    return (
        <div className={styles.overlay} onClick={onClose}>
            <div className={styles.modal} onClick={e => e.stopPropagation()}>
                <h2 className={styles.title}>Вышло обновление v{version}</h2>
                <p className={styles.text}>В приложение добавлены новые функции. Посмотрите, что изменилось.</p>
                <div className={styles.actions}>
                    <button className={styles.primary} onClick={() => { onClose(); navigate('/changelog'); }}>
                        Посмотреть
                    </button>
                    <button className={styles.secondary} onClick={onClose}>
                        Закрыть
                    </button>
                </div>
            </div>
        </div>
    );
}
