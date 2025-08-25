import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LoginPage from '../pages/LoginPage'
import RegisterPage from '../pages/RegisterPage'
import HomePage from '../pages/HomePage'
import NoteList from '../pages/NoteList'
import RequireAuth from '../features/auth/components/RequireAuth'
import NoteDetailPage from '../pages/NoteDetailPage'
import Layout from './Layout'

const AppRouter = () => (
    <BrowserRouter>
        <Routes>
            {/* публичные страницы */}
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />

            {/* защищённые страницы */}
            <Route
                path="/"
                element={
                    <RequireAuth>
                        <Layout />
                    </RequireAuth>
                }
            >
                <Route index element={<HomePage />} />
                <Route path="notes" element={<NoteList />} />
                <Route path="notes/:id" element={<NoteDetailPage />} />
                <Route path="archive" element={<NoteList isArchiveView />} />
            </Route>
        </Routes>
    </BrowserRouter>
)

export default AppRouter
