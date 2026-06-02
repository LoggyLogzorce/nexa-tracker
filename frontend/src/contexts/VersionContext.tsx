import { createContext, useState, useEffect, useCallback, type ReactNode } from 'react';
import { getVersion } from '../api/version';

interface VersionContextType {
    hasNewVersion: boolean;
    currentVersion: string;
    markSeen: () => void;
}

const LAST_SEEN_KEY = 'last_seen_version';

// eslint-disable-next-line react-refresh/only-export-components
export const VersionContext = createContext<VersionContextType | null>(null);

export function VersionProvider({ children }: { children: ReactNode }) {
    const [currentVersion, setCurrentVersion] = useState('');
    const [hasNewVersion, setHasNewVersion] = useState(false);

    useEffect(() => {
        let cancelled = false;
        getVersion().then(v => {
            if (cancelled) return;
            setCurrentVersion(v);
            const lastSeen = localStorage.getItem(LAST_SEEN_KEY);
            if (lastSeen !== v) {
                setHasNewVersion(true);
            }
        }).catch(() => {});
        return () => { cancelled = true; };
    }, []);

    const markSeen = useCallback(() => {
        localStorage.setItem(LAST_SEEN_KEY, currentVersion);
        setHasNewVersion(false);
    }, [currentVersion]);

    return (
        <VersionContext.Provider value={{ hasNewVersion, currentVersion, markSeen }}>
            {children}
        </VersionContext.Provider>
    );
}
