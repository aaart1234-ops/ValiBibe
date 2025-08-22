import React from 'react'
import {
    AppBar,
    Box,
    Toolbar,
    Typography,
    IconButton,
    Drawer,
    List,
    ListItemButton,
    ListItemText
} from '@mui/material'
import MenuIcon from '@mui/icons-material/Menu'
import { Link as RouterLink, useNavigate } from 'react-router-dom'
import logo from '../assets/logo.png'

const Header = () => {
    const [drawerOpen, setDrawerOpen] = React.useState(false)
    const navigate = useNavigate()

    const toggleDrawer = (open: boolean) => () => {
        setDrawerOpen(open)
    }

    const handleNavigation = (path: string) => () => {
        setDrawerOpen(false)
        navigate(path)
    }

    return (
        <>
            <AppBar position="static">
                <Toolbar>
                    {/* Логотип */}
                    <Typography
                        variant="h6"
                        component={RouterLink}
                        to="/"
                        sx={{
                            flexGrow: 1,
                            textDecoration: 'none',
                            color: 'inherit'
                        }}
                    >
                        {/*<img
                            src={logo}
                            alt="Logo"
                            style={{height: 30, marginRight: 8, marginTop: 10}} // регулируй размер и отступ
                        />*/}
                        <b>ValiBibe</b>
                    </Typography>

                    {/* Иконка бургера */}
                    <IconButton
                        edge="end"
                        color="inherit"
                        onClick={toggleDrawer(true)}
                    >
                        <MenuIcon />
                    </IconButton>
                </Toolbar>
            </AppBar>

            {/* Drawer-меню */}
            <Drawer anchor="right" open={drawerOpen} onClose={toggleDrawer(false)}>
                <Box
                    sx={{ width: 250 }}
                    role="presentation"
                    onClick={toggleDrawer(false)}
                    onKeyDown={toggleDrawer(false)}
                >
                    <List>
                        <ListItemButton onClick={handleNavigation('/')}>
                            <ListItemText primary="Главная" />
                        </ListItemButton>
                        <ListItemButton onClick={handleNavigation('/notes')}>
                            <ListItemText primary="Заметки" />
                        </ListItemButton>
                        <ListItemButton onClick={handleNavigation('/archive')}>
                            <ListItemText primary="Показать архив" />
                        </ListItemButton>
                        <ListItemButton onClick={handleNavigation('/login')}>
                            <ListItemText primary="Выйти" />
                        </ListItemButton>
                    </List>
                </Box>
            </Drawer>
        </>
    )
}

export default Header
