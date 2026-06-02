import client from './client';

export const getVersion = async (): Promise<string> => {
    const res = await client.get<{ version: string }>('/version');
    return res.data.version;
};

export const getChangelog = async (): Promise<string> => {
    const res = await client.get('/changelog');
    return res.data;
};
