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
const user = localStorage.getItem('user')

if (token && user) {
    store.dispatch(setCredentials({ token, user: JSON.parse(user) }))
}

// Типы для useDispatch и useSelector
export type RootState = ReturnType<typeof store.getState>
export type AppDispatch = typeof store.dispatch
