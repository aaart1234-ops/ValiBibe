// features/note/components/SwipeableNoteCard.tsx
import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import {
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
    Box,
    Menu,
    MenuItem,
} from '@mui/material'
import {
    SwipeableList,
    SwipeableListItem,
    LeadingActions,
    TrailingActions,
    SwipeAction,
} from 'react-swipeable-list'
import 'react-swipeable-list/dist/styles.css'

import { Note, useDeleteNoteMutation, useUnarchiveNoteMutation } from '../noteApi'
import NoteCard from './NoteCard'
import { useLongPress } from '@/features/note/hooks/useLongPress'

interface SwipeableNoteCardProps {
    note: Note
    onRefetch?: () => void
    onRequestArchive?: (note: Note) => void
}

const SwipeableNoteCard: React.FC<SwipeableNoteCardProps> = ({ note, onRefetch, onRequestArchive }) => {
    const [deleteNote] = useDeleteNoteMutation()
    const [unarchiveNote] = useUnarchiveNoteMutation()

    const [confirmDelete, setConfirmDelete] = useState(false)
    const [isProcessing, setIsProcessing] = useState(false)

    // --- состояние для long press меню
    const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null)

    const handleUnarchive = async () => {
        try {
            await unarchiveNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при разархивировании:', e)
        }
    }

    const handleDelete = async () => {
        if (isProcessing) return
        setIsProcessing(true)
        try {
            await deleteNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при удалении:', e)
        } finally {
            setConfirmDelete(false)
            setIsProcessing(false)
        }
    }

    // --- long press
    const longPressHandlers = useLongPress({
        onLongPress: (_e, target) => {
            setMenuAnchor(target)   // теперь точно будет HTMLElement
        },
        delay: 600,
    })

    return (
        <>
            <SwipeableList threshold={0.25}>
                <SwipeableListItem
                    leadingActions={
                        <LeadingActions>
                            <SwipeAction
                                onClick={() =>
                                    note.archived ? handleUnarchive() : onRequestArchive?.(note)
                                }
                            >
                                <Box
                                    display="flex"
                                    justifyContent="center"
                                    alignItems="center"
                                    bgcolor="rgb(175, 238, 175)"
                                    width="100%"
                                    height="100%"
                                >
                                    {note.archived ? 'Из архива' : 'В архив'}
                                </Box>
                            </SwipeAction>
                        </LeadingActions>
                    }
                    trailingActions={
                        <TrailingActions>
                            <SwipeAction onClick={() => setConfirmDelete(true)}>
                                <Box
                                    display="flex"
                                    justifyContent="center"
                                    alignItems="center"
                                    bgcolor="rgb(238, 175, 175)"
                                    width="100%"
                                    height="100%"
                                >
                                    Удалить
                                </Box>
                            </SwipeAction>
                        </TrailingActions>
                    }
                    blockSwipe={confirmDelete}
                >
                    {/* Оборачиваем Link, чтобы ловить long press */}
                    <Box
                        {...longPressHandlers}
                        onContextMenu={(e) => e.preventDefault()}
                        sx={{ flexGrow: 1, display: 'flex' }}
                    >
                        <Link
                            to={`/notes/${note.id}`}
                            style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}
                        >
                            <NoteCard note={note} />
                        </Link>
                    </Box>
                </SwipeableListItem>
            </SwipeableList>

            {/* Контекстное меню по long press */}
            <Menu
                anchorEl={menuAnchor}
                open={Boolean(menuAnchor)}
                onClose={() => setMenuAnchor(null)}
            >
                <MenuItem
                    onClick={() => {
                        setMenuAnchor(null)
                        note.archived ? handleUnarchive() : onRequestArchive?.(note)
                    }}
                >
                    {note.archived ? 'Из архива' : 'В архив'}
                </MenuItem>
                <MenuItem
                    onClick={() => {
                        setMenuAnchor(null)
                        setConfirmDelete(true)
                    }}
                >
                    Удалить
                </MenuItem>
            </Menu>

            {/* Диалог удаления */}
            <Dialog open={confirmDelete} onClose={() => setConfirmDelete(false)}>
                <DialogTitle>Удалить заметку?</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Вы уверены, что хотите удалить эту заметку? Это действие необратимо.
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDelete(false)} disabled={isProcessing}>
                        Отмена
                    </Button>
                    <Button onClick={handleDelete} color="error" disabled={isProcessing}>
                        Удалить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )
}

export default SwipeableNoteCard
