import { useState } from 'react'
import {
    Card,
    CardContent,
    Typography,
    Box,
    LinearProgress,
    Chip,
    Menu,
    MenuItem,
} from '@mui/material'
import { Note } from '../../noteApi'
import dayjs from 'dayjs'
import EventIcon from '@mui/icons-material/Event'
import ArchiveIcon from '@mui/icons-material/Archive'
import { useLongPress } from '@/features/note/hooks/useLongPress'

interface NoteCardProps {
    note: Note
    onRequestArchive?: (note: Note) => void // опционально
    onRequestDelete?: (note: Note) => void // опционально
    isArchiveView?: boolean
}

const NoteCard = ({
                      note,
                      onRequestArchive,
                      onRequestDelete,
                      isArchiveView,
                  }: NoteCardProps) => {
    const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

    const handleOpenMenu = (e: any) => setAnchorEl(e.currentTarget as HTMLElement)
    const handleClose = () => setAnchorEl(null)

    const enableMenu = Boolean(onRequestArchive || onRequestDelete)

    const longPressHandlers = enableMenu
        ? useLongPress({
            onLongPress: handleOpenMenu,
            delay: 600,
        })
        : ({} as Record<string, unknown>)

    return (
        <Card
            sx={{ display: 'flex', flexDirection: 'column', flexGrow: 1, position: 'relative' }}
            {...longPressHandlers}
        >
            {/* Лейбл "В архиве" */}
            {note.archived && (
                <Chip
                    icon={<ArchiveIcon />}
                    label="В архиве"
                    size="small"
                    color="default"
                    sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1 }}
                    onClick={(e) => e.stopPropagation()}
                />
            )}

            <CardContent sx={{ flexGrow: 1, display: 'flex', flexDirection: 'column' }}>
                <Typography
                    variant="h6"
                    sx={{
                        display: '-webkit-box',
                        WebkitLineClamp: 2,
                        WebkitBoxOrient: 'vertical',
                        overflow: 'hidden',
                        textOverflow: 'ellipsis',
                    }}
                >
                    {note.title}
                </Typography>

                <Box display="flex" alignItems="center" gap={0.5} mt={0.5} mb={4}>
                    <EventIcon sx={{ fontSize: 18, color: 'text.secondary' }} />
                    <Typography variant="body2" color="text.secondary">
                        {dayjs(note.created_at).format('D MMMM YYYY')}
                    </Typography>
                </Box>

                <Box mt="auto">
                    <Typography variant="body2" color="text.secondary">
                        🧠 Уровень запоминания: {note.memoryLevel}%
                    </Typography>
                    <LinearProgress
                        variant="determinate"
                        value={note.memoryLevel}
                        sx={{ height: 8, borderRadius: 5, mt: 0.5 }}
                        color={note.memoryLevel < 40 ? 'error' : note.memoryLevel < 70 ? 'warning' : 'success'}
                    />
                </Box>
            </CardContent>

            {/* Меню действий (только если прокинуты коллбеки) */}
            <Menu anchorEl={anchorEl} open={Boolean(anchorEl) && enableMenu} onClose={handleClose}>
                {onRequestDelete && (
                    <MenuItem
                        onClick={() => {
                            handleClose()
                            onRequestDelete(note)
                        }}
                    >
                        Удалить
                    </MenuItem>
                )}
                {onRequestArchive && (
                    <MenuItem
                        onClick={() => {
                            handleClose()
                            onRequestArchive(note)
                        }}
                    >
                        {isArchiveView ? 'Вернуть из архива' : 'Архивировать'}
                    </MenuItem>
                )}
            </Menu>
        </Card>
    )
}

export default NoteCard
