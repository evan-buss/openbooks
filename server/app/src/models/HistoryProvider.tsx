import useLocalStorage from '@rehooks/local-storage';
import React from 'react';
import { BookDetail } from './messages';


// HistoryItem represents a single search history item
export type HistoryItem = {
    query: string;
    time: number;
    results?: BookDetail[];
}

// HistoryContextType allows accessing current history and adding new items
type HistoryContextType = {
    history: HistoryItem[];
    addItem: (value: HistoryItem) => void;
}


export const HistoryContext = React.createContext<HistoryContextType | undefined>(undefined);
// useHistory is a helper hook to access HistoryContext
export const useHistory = () => React.useContext(HistoryContext);

type Props = {
    children: React.ReactNode
}

// HOC that wraps HistoryContextProvider
const HistoryProvider: React.FC<Props> = ({ children }: Props) => {

    const [history, setHistory] = useLocalStorage<HistoryItem[]>("history", []);

    // Each time a new item is added, remove the oldest item if over 15 items
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