// src/features/auth/authSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit'
import { noteApi } from '../note/noteApi' // важно
import { AppDispatch } from '@/app/store'  // импорт типа для thunk

interface AuthState {
    token: string | null
}

const initialState: AuthState = {
    token: null,
}

const authSlice = createSlice({
    name: 'auth',
    initialState,
    reducers: {
        setCredentials: (state, action: PayloadAction<{ token: string }>) => {
            state.token = action.payload.token
        },
        clearCredentials: (state) => {
            state.token = null
        },
    },
})

// 👇 Thunk-действие для logout с очисткой кэша
export const logout = () => (dispatch: AppDispatch) => {
    dispatch(authSlice.actions.clearCredentials())
    dispatch(noteApi.util.resetApiState())
}

export const { setCredentials } = authSlice.actions
export default authSlice.reducer
