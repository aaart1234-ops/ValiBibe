import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import {
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

import { Note, useUnarchiveNoteMutation } from '../../noteApi'
import NoteRow from './NoteRow'
import { useLongPress } from '../../hooks/useLongPress'

interface SwipeableNoteRowProps {
    note: Note
    onRefetch?: () => void
    onRequestArchive?: (note: Note) => void
    onRequestDelete?: (note: Note) => void
}

const SwipeableNoteRow: React.FC<SwipeableNoteRowProps> = ({
                                                               note,
                                                               onRefetch,
                                                               onRequestArchive,
                                                               onRequestDelete,
                                                           }) => {
    const [unarchiveNote] = useUnarchiveNoteMutation()

    const [menuAnchor, setMenuAnchor] = useState<null | HTMLElement>(null)

    const handleUnarchive = async () => {
        try {
            await unarchiveNote(note.id).unwrap()
            onRefetch?.()
        } catch (e) {
            console.error('Ошибка при разархивировании:', e)
        }
    }

    // long-press хэндлеры
    const longPressHandlers = useLongPress({
        onLongPress: (_e, target) => {
            setMenuAnchor(target)
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
                            <SwipeAction onClick={() => onRequestDelete?.(note)}>
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
                >
                    <Box
                        {...longPressHandlers}
                        onContextMenu={(e) => e.preventDefault()}
                        sx={{ display: 'block', width: '100%' }}
                    >
                        <Link
                            to={`/notes/${note.id}`}
                            style={{ textDecoration: 'none', display: 'block', width: '100%' }}
                        >
                            <NoteRow note={note} />
                        </Link>
                    </Box>
                </SwipeableListItem>
            </SwipeableList>

            {/* Меню по долгому тапу */}
            <Menu
                anchorEl={menuAnchor}
                open={Boolean(menuAnchor)}
                onClose={() => setMenuAnchor(null)}
            >
                {note.archived ? (
                    <MenuItem
                        onClick={() => {
                            handleUnarchive()
                            setMenuAnchor(null)
                        }}
                    >
                        Разархивировать
                    </MenuItem>
                ) : (
                    <MenuItem
                        onClick={() => {
                            onRequestArchive?.(note)
                            setMenuAnchor(null)
                        }}
                    >
                        В архив
                    </MenuItem>
                )}
                <MenuItem
                    onClick={() => {
                        onRequestDelete?.(note)
                        setMenuAnchor(null)
                    }}
                >
                    Удалить
                </MenuItem>
            </Menu>
        </>
    )
}

export default SwipeableNoteRow
