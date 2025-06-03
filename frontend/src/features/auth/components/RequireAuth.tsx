import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'

const RequireAuth = ({ children }: { children: React.ReactNode }) => {
    const token = localStorage.getItem('token')
    const location = useLocation()

    if (!token) {
        // если не авторизован — редиректим на /login
        return <Navigate to="/login" state={{ from: location }} replace />
    }

    return children
}

export default RequireAuth
