import { useContext } from 'react';
import { VersionContext } from './VersionContext';

export function useVersion() {
    const ctx = useContext(VersionContext);
    if (!ctx) throw new Error('useVersion must be used within VersionProvider');
    return ctx;
}
