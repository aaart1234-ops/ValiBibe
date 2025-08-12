// features/note/components/SwipeableNoteRow.tsx
import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import {
    Grid,
    Dialog,
    DialogTitle,
    DialogContent,
    DialogContentText,
    DialogActions,
    Button,
} from '@mui/material'
import {
    SwipeableList,
    SwipeableListItem,
    LeadingActions,
    TrailingActions,
    SwipeAction,
} from 'react-swipeable-list'
import 'react-swipeable-list/dist/styles.css'

import { Note, useDeleteNoteMutation, useArchiveNoteMutation, useUnarchiveNoteMutation } from '../noteApi'
import NoteRow from './NoteRow'

interface SwipeableNoteRowProps {
    note: Note
    onRefetch?: () => void
    onRequestArchive?: (note: Note) => void
}

const SwipeableNoteRow: React.FC<SwipeableNoteRowProps> = ({ note, onRefetch, onRequestArchive }) => {
    const [deleteNote] = useDeleteNoteMutation()
    const [unarchiveNote] = useUnarchiveNoteMutation()
    // archive is delegated to parent (onRequestArchive) for optimistic + undo
    // const [archiveNote] = useArchiveNoteMutation() // not used here

    const [confirmDialog, setConfirmDialog] = useState<'archive' | 'delete' | null>(null)
    const [isProcessing, setIsProcessing] = useState(false)

    const handleConfirm = async () => {
        setIsProcessing(true)
        try {
            if (confirmDialog === 'delete') {
                await deleteNote(note.id).unwrap()
            } else if (confirmDialog === 'archive') {
                // delegate to parent if provided
                onRequestArchive?.(note)
            }
            onRefetch?.()
        } catch (error) {
            console.error('Ошибка:', error)
        } finally {
            setConfirmDialog(null)
            setIsProcessing(false)
        }
    }

    const handleArchiveToggle = () => {
        if (note.archived) {
            // разархивирование без диалога
            unarchiveNote(note.id)
                .unwrap()
                .then(() => onRefetch?.())
                .catch((err) => console.error('Ошибка разархивирования:', err))
        } else {
            // архивирование — делегируем родителю
            onRequestArchive?.(note)
        }
    }

    return (
        <>
            <SwipeableList threshold={0.25}>
                <SwipeableListItem
                    leadingActions={
                        <LeadingActions>
                            <SwipeAction onClick={handleArchiveToggle}>
                                <Grid
                                    container
                                    key={note.id}
                                    sx={{
                                        width: { xs: '100%', sm: '48%', md: '32%' },
                                        backgroundColor: 'rgb(175, 238, 175)',
                                    }}
                                    display="flex"
                                    alignItems="center"
                                    justifyContent="center"
                                >
                                    {note.archived ? 'Извлечь из архива' : 'В архив'}
                                </Grid>
                            </SwipeAction>
                        </LeadingActions>
                    }
                    trailingActions={
                        <TrailingActions>
                            <SwipeAction onClick={() => setConfirmDialog('delete')}>
                                <Grid
                                    container
                                    key={note.id}
                                    sx={{
                                        width: { xs: '100%', sm: '48%', md: '32%' },
                                        backgroundColor: 'rgb(238, 175, 175)',
                                    }}
                                    display="flex"
                                    alignItems="center"
                                    justifyContent="center"
                                >
                                    Удалить
                                </Grid>
                            </SwipeAction>
                        </TrailingActions>
                    }
                >
                    <Link key={note.id} to={`/notes/${note.id}`} style={{ textDecoration: 'none', flexGrow: 1, display: 'flex' }}>
                        <NoteRow note={note} />
                    </Link>
                </SwipeableListItem>
            </SwipeableList>

            <Dialog open={!!confirmDialog} onClose={() => setConfirmDialog(null)}>
                <DialogTitle>
                    {confirmDialog === 'delete' ? 'Удалить заметку?' : 'Архивировать заметку?'}
                </DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        {confirmDialog === 'delete'
                            ? 'Вы уверены, что хотите удалить эту заметку? Это действие необратимо.'
                            : 'После архивирования заметка будет скрыта из основного списка.'}
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setConfirmDialog(null)} disabled={isProcessing}>Отмена</Button>
                    <Button onClick={handleConfirm} color={confirmDialog === 'delete' ? 'error' : 'primary'} disabled={isProcessing}>
                        Подтвердить
                    </Button>
                </DialogActions>
            </Dialog>
        </>
    )
}

export default SwipeableNoteRow
