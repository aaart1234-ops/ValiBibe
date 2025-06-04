import { useState, useEffect } from 'react'
import { TextField, Button, Box, Typography, Alert } from '@mui/material'
import { useLoginMutation } from '../authApi'
import { useAppDispatch } from '@/app/hooks'
import { setCredentials } from '../authSlice'
import { useNavigate, useLocation, Link } from 'react-router-dom'

const LoginForm = () => {
    const dispatch = useAppDispatch()
    const navigate = useNavigate()
    const location = useLocation()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [successMessage, setSuccessMessage] = useState<string | null>(null)

    const [login, { isLoading, error }] = useLoginMutation()

    // ✅ Отображаем successMessage, переданное из /register
    useEffect(() => {
        if (location.state?.successMessage) {
            setSuccessMessage(location.state.successMessage)

            // Удаляем message из history, чтобы не оставалось при обновлении страницы
            window.history.replaceState({}, document.title)
        }
    }, [location.state])

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            const response = await login({ email, password }).unwrap()
            dispatch(setCredentials(response))
            localStorage.setItem('token', response.token)
            navigate('/')
        } catch (err) {
            console.error('Login failed:', err)
        }
    }

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 400, mx: 'auto', mt: 6 }}>
            <Typography variant="h5" mb={2}>
                Вход в аккаунт
            </Typography>

            {successMessage && (
                <Alert severity="success" sx={{ mb: 2 }}>
                    {successMessage}
                </Alert>
            )}

            <TextField
                label="Email"
                type="email"
                fullWidth
                required
                margin="normal"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
            />

            <TextField
                label="Пароль"
                type="password"
                fullWidth
                required
                margin="normal"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
            />

            {error && (
                <Alert severity="error" sx={{ mt: 2 }}>
                    Неверный email или пароль
                </Alert>
            )}

            <Button
                type="submit"
                variant="contained"
                color="primary"
                fullWidth
                sx={{ mt: 2 }}
                disabled={isLoading}
            >
                {isLoading ? 'Вход...' : 'Войти'}
            </Button>

            <Typography variant="body2" align="center" sx={{ mt: 2 }}>
                Ещё нет аккаунта?{' '}
                <Link to="/register">Зарегистрироваться</Link>
            </Typography>
        </Box>
    )
}

export default LoginForm
