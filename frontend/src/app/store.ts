import { configureStore } from '@reduxjs/toolkit'
import authReducer from '../features/auth/authSlice'
import { authApi } from '../features/auth/authApi'
import { setCredentials } from '../features/auth/authSlice'

export const store = configureStore({
    reducer: {
        auth: authReducer,
        [authApi.reducerPath]: authApi.reducer,
    },
    middleware: (getDefaultMiddleware) =>
        getDefaultMiddleware().concat(authApi.middleware),
})

// Попытка восстановить сессию
const token = localStorage.getItem('token')

if (token) {
    try {
        store.dispatch(setCredentials({ token }))
    } catch (e) {
        console.error('Ошибка при чтении user из localStorage:', e)
        // Очистим некорректные данные
        localStorage.removeItem('token')
    }
}

// Типы для useDispatch и useSelector
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
