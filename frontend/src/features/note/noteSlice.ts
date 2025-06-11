// features/note/noteSlice.ts
import { createSlice, PayloadAction } from '@reduxjs/toolkit'

const noteSlice = createSlice({
    name: 'note',
    initialState: {
        viewMode: 'card' as 'card' | 'row',
    },
    reducers: {
        toggleViewMode: (state) => {
            state.viewMode = state.viewMode === 'card' ? 'row' : 'card'
        },
    },
})

export const { toggleViewMode } = noteSlice.actions
export default noteSlice.reducer
