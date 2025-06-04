import {useState} from 'react'
import {TextField, Button, Box, Typography, Alert} from '@mui/material'
import { useRegisterMutation, useLoginMutation } from '../authApi'
import { useAppDispatch } from '@/app/hooks'
import { setCredentials } from '../authSlice'
import { useNavigate, Link } from 'react-router-dom'

const RegisterForm = () => {
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [nickname, setNickname] = useState('')

    const [emailError, setEmailError] = useState('')
    const [passwordError, setPasswordError] = useState('')
    const [formError, setFormError] = useState('')

    const navigate = useNavigate()

    //Добавим вызов и состояния
    const [register, {isLoading}] = useRegisterMutation()
    const [login] = useLoginMutation()
    const dispatch = useAppDispatch()

    const validateEmail = (value: string) => {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
        return emailRegex.test(value)
    }

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        // Сброс ошибок
        setEmailError('')
        setPasswordError('')
        setFormError('')

        let valid = true

        if (!validateEmail(email)) {
            setEmailError('Введите корректный email')
            valid = false
        }

        if (password.length < 6) {
            setPasswordError('Пароль должен содержать минимум 6 символов')
            valid = false
        }

        if (!valid) return

        //Логика отправки запроса
        try {
            // 1. Регистрируем пользователя
            const res = await register({email, password, nickname}).unwrap()

            // 2. Логиним сразу же
            const loginRes = await login({ email, password }).unwrap()

            dispatch(setCredentials(loginRes))
            localStorage.setItem('token', loginRes.token)

            navigate('/')

            // Здесь можешь делать redirect или сброс формы
        } catch (err: any) {
            const message = err?.data?.error || 'Произошла ошибка регистрации'
            setFormError(message)
        }
    }

    return (
        <Box component="form" onSubmit={handleSubmit} sx={{ maxWidth: 400, mx: 'auto', mt: 6 }}>
            <Typography variant="h5" mb={2}>
                Регистрация
            </Typography>

            <TextField
                label="Имя пользователя"
                type="text"
                fullWidth
                //required
                margin="normal"
                value={nickname}
                onChange={(e) => setNickname(e.target.value)}
            />

            <TextField
                label="Email"
                type="email"
                fullWidth
                //required
                margin="normal"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                error={!!emailError}
                helperText={emailError}
            />

            <TextField
                label="Пароль"
                type="password"
                fullWidth
                required
                margin="normal"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                error={!!passwordError}
                helperText={passwordError}
            />

            {formError && (
                <Alert severity="error" sx={{ mt: 2 }}>
                    {formError}
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
                Зарегистрироваться
            </Button>

            <Typography variant="body2" align="center" sx={{ mt: 2 }}>
                Уже есть аккаунт?{' '}
                <Link to="/login">Войти</Link>
            </Typography>
        </Box>
    )
}

export default RegisterForm