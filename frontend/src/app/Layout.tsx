import { Outlet } from 'react-router-dom'
import Header from '@/components/Header'
import { DetailPageProvider } from '@/context/DetailPageContext'

export default function Layout() {
    return (
        <DetailPageProvider>
            <Header />
            <Outlet />
        </DetailPageProvider>
    )
}
