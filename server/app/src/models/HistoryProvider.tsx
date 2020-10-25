import useLocalStorage from '@rehooks/local-storage';
import React from 'react';
import { BookDetail } from './messages';


// HistoryItem
export type HistoryItem = {
    query: string;
    time: number;
    results?: BookDetail[];
}

type HistoryContextType = {
    history: HistoryItem[];
    addItem: (value: HistoryItem) => void;
}

export const HistoryContext = React.createContext<HistoryContextType | undefined>(undefined);
export const useHistory = () => React.useContext(HistoryContext);

type Props = {
    children: React.ReactNode
}

const HistoryProvider: React.FC<Props> = ({ children }: Props) => {

    const [history, setHistory] = useLocalStorage<HistoryItem[]>("history", []);

    const addItem = (newItem: HistoryItem) => {
        setHistory([newItem, ...history].slice(0, 16));
    }

    return (
        <HistoryContext.Provider value={{ history, addItem }}>
            { children}
        </ HistoryContext.Provider>
    )
}

export default HistoryProvider;