import { BrowserRouter, Routes, Route } from 'react-router-dom'
import LoginPage from '../pages/LoginPage'

const AppRouter = () => (
    <BrowserRouter>
        <Routes>
            <Route path="/login" element={<LoginPage />} />
            {/* другие страницы добавим позже */}
        </Routes>
    </BrowserRouter>
)

export default AppRouter
