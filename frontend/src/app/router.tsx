import { BrowserRouter, Routes, Route, useNavigate  } from 'react-router-dom'
import { useEffect } from 'react'
import LoginPage from '../pages/LoginPage'
import RegisterPage from '../pages/RegisterPage'
import HomePage from '../pages/HomePage'
import NoteList from '../pages/NoteList'
import RequireAuth from '../features/auth/components/RequireAuth'
import Header from '@/components/Header'
import NoteDetailPage from '../pages/NoteDetailPage'

const AppRouter = () => (
    <BrowserRouter>
        <Header />
        <Routes>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
            <Route
                path="/"
                element={
                    <RequireAuth>
                        <HomePage />
                    </RequireAuth>
                }
            />
            <Route
                path="/notes"
                element={
                    <RequireAuth>
                        <NoteList />
                    </RequireAuth>
                }
            />
            <Route
                path="/notes/:id"
                element={
                    <RequireAuth>
                        <NoteDetailPage />
                    </RequireAuth>
                }
            />
            <Route
                path="/archive"
                element={
                    <RequireAuth>
                        <NoteList isArchiveView />
                    </RequireAuth>
                }
            />
            {/* другие страницы */}
        </Routes>
    </BrowserRouter>
)

export default AppRouter
