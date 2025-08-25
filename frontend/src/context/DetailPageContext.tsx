import React, { createContext, useContext, useMemo, useState } from 'react'

type DetailPageCtx = {
    isEditing: boolean
    setEditing: (v: boolean) => void
}

const DetailPageContext = createContext<DetailPageCtx | undefined>(undefined)

export const DetailPageProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
    const [isEditing, setEditing] = useState(false)

    const value = useMemo(() => ({ isEditing, setEditing }), [isEditing])
    return <DetailPageContext.Provider value={value}>{children}</DetailPageContext.Provider>
}

export const useDetailPage = () => {
    const ctx = useContext(DetailPageContext)
    if (!ctx) {
        throw new Error('useDetailPage must be used within DetailPageProvider')
    }
    return ctx
}
