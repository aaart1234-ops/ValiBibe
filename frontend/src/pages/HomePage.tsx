// src/pages/HomePage.tsx

import React from 'react'
import { Box, Typography, Button } from '@mui/material'
import { useNavigate } from 'react-router-dom'

const HomePage = () => {
    const navigate = useNavigate()

    const handleLogout = () => {
        localStorage.removeItem('token')
        localStorage.removeItem('user')
        navigate('/login')
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
