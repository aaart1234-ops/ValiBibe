import React from 'react'
import {
    AppBar, Box, Toolbar, Typography, IconButton, Drawer, List, ListItemButton, ListItemText,
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import { Link as RouterLink, useNavigate, useLocation } from 'react-router-dom'
import ArrowBackIcon from '@mui/icons-material/ArrowBack'
import { useDetailPage } from '@/context/DetailPageContext'

const Header = () => {
    const [drawerOpen, setDrawerOpen] = React.useState(false)
    const navigate = useNavigate()
    const location = useLocation()
    const { isEditing, setEditing } = useDetailPage()

    const isDetailPage = location.pathname.startsWith('/notes/')

    const handleBack = () => {
        if (isDetailPage && isEditing) {
            // вместо навигации — выходим из режима редактирования
            setEditing(false)
        } else {
            navigate(-1)
        }
    }

    return (
        <>
            <AppBar position="sticky">
                <Toolbar>
                    {isDetailPage ? (
                        <IconButton
                            edge="start"
                            color="inherit"
                            onClick={handleBack}
                            sx={{ flexGrow: 1, display: 'flex', justifyContent: 'flex-start' }}
                        >
                            <ArrowBackIcon />
                        </IconButton>
                    ) : (
                        <Typography
                            variant="h6"
                            component={RouterLink}
                            to="/notes"
                            sx={{ flexGrow: 1, textDecoration: 'none', color: 'inherit' }}
                        >
                            <b>V</b>
                        </Typography>
                    )}
                    <IconButton edge="end" color="inherit" onClick={() => setDrawerOpen(true)}>
                        <MenuIcon />
                    </IconButton>
                </Toolbar>
            </AppBar>

            <Drawer anchor="right" open={drawerOpen} onClose={() => setDrawerOpen(false)}>
                <Box sx={{ width: 250 }} role="presentation">
                    <List>
                        <ListItemButton onClick={() => navigate('/')}>
                            <ListItemText primary="Главная" />
                        </ListItemButton>
                        <ListItemButton onClick={() => navigate('/notes')}>
                            <ListItemText primary="Заметки" />
                        </ListItemButton>
                        <ListItemButton onClick={() => navigate('/archive')}>
                            <ListItemText primary="Показать архив" />
                        </ListItemButton>
                        <ListItemButton onClick={() => navigate('/login')}>
                            <ListItemText primary="Выйти" />
                        </ListItemButton>
                    </List>
                </Box>
            </Drawer>
        </>
    )
}

export default Header
