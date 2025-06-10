// src/features/auth/components/RequireAuth.tsx
import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAppSelector } from '@/app/hooks'

const RequireAuth = ({ children }: { children: React.ReactNode }) => {
    const token = useAppSelector((state) => state.auth.token)
    const location = useLocation()

    if (!token) {
        // если не авторизован — редиректим на /login
        return <Navigate to="/login" state={{ from: location }} replace />
    }

    return <>{children}</>
}

export default RequireAuth
