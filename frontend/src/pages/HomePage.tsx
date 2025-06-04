// src/pages/HomePage.tsx

import React from 'react'
import { Box, Typography, Button } from '@mui/material'
import { useNavigate } from 'react-router-dom'
import { useAppDispatch } from '@/app/hooks'
import {logout} from "@/features/auth/authSlice";


const HomePage = () => {
    const navigate = useNavigate()
    const dispatch = useAppDispatch()

    const handleLogout = () => {
        dispatch(logout())  // <--- очищаем Redux-состояние
        localStorage.removeItem('token')  // удаляем токен из localStorage
        navigate('/login') // редирект на логин
    }

    return (
        <Box sx={{ maxWidth: 600, mx: 'auto', mt: 8, textAlign: 'center' }}>
            <Typography variant="h4" gutterBottom>
                Добро пожаловать!
            </Typography>
            <Typography variant="body1" sx={{ mb: 4 }}>
                Вы успешно вошли в систему.
            </Typography>
            <Button variant="contained" color="primary" onClick={handleLogout}>
                Выйти
            </Button>
        </Box>
    )
}

export default HomePage
