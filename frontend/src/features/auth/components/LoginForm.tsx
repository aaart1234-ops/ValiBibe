import { useState } from 'react'
import { TextField, Button, Box, Typography, Alert } from '@mui/material'
import { useLoginMutation } from '../authApi'
import { useAppDispatch } from '@/app/hooks'
import { setCredentials } from '../authSlice'
import { useNavigate } from 'react-router-dom'

const LoginForm = () => {
    const dispatch = useAppDispatch()
    const navigate = useNavigate()

    // Поля формы
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')

    // Ошибки валидации
    const [emailError, setEmailError] = useState('')
    const [passwordError, setPasswordError] = useState('')

    const [login, { isLoading, error }] = useLoginMutation()

    // Валидация email
    const validateEmail = () => {
        if (!email) {
            setEmailError('Введите email')
            return false
        }
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        if (!emailRegex.test(email)) {
            setEmailError('Некорректный email')
            return false
        }
        setEmailError('')
        return true
    }

    // Валидация пароля
    const validatePassword = () => {
        if (!password) {
            setPasswordError('Введите пароль')
            return false
        }
        if (password.length < 6) {
            setPasswordError('Пароль должен быть не менее 6 символов')
            return false
        }
        setPasswordError('')
        return true
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // Запускаем валидацию
        const isEmailValid = validateEmail()
        const isPasswordValid = validatePassword()

        if (!isEmailValid || !isPasswordValid) return

        try {
            const response = await login({ email, password }).unwrap()
            // 1. сохраняем в Redux
            dispatch(setCredentials(response))

            // 2. сохраняем в localStorage
            localStorage.setItem('token', response.token)
            localStorage.setItem('user', JSON.stringify(response.user))

            // 3. редирект на список заметок
            navigate('/notes')
        } catch (err) {
            console.error('Login failed:', err)
        }
    }

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 400, mx: 'auto', mt: 6 }}>
            <Typography variant="h5" mb={2}>
                Вход в аккаунт
            </Typography>

            <TextField
                label="Email"
                type="email"
                fullWidth
                margin="normal"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                onBlur={validateEmail} // ← Валидация при потере фокуса
                error={!!emailError}
                helperText={emailError}
            />

            <TextField
                label="Пароль"
                type="password"
                fullWidth
                margin="normal"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                onBlur={validatePassword} // ← Валидация при потере фокуса
                error={!!passwordError}
                helperText={passwordError}
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
        </Box>
    )
}

export default LoginForm
