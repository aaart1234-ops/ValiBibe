import { Box, Typography, TextField, IconButton } from '@mui/material'
import EditIcon from '@mui/icons-material/Edit'
import MoreVertIcon from '@mui/icons-material/MoreVert'
import { RefObject } from 'react'

interface NoteTitleProps {
    title: string
    isEditing: boolean
    onChange: (value: string) => void
    onEditToggle: (val: boolean) => void
    titleRef: RefObject<HTMLInputElement>
    onMoreClick: (e: React.MouseEvent<HTMLButtonElement>) => void // ✅ новый проп
}

export const NoteTitle = ({
                              title,
                              isEditing,
                              onChange,
                              onEditToggle,
                              titleRef,
                              onMoreClick,
                          }: NoteTitleProps) => {
    return (
        <Box>
            {!isEditing ? (
                <Box display="flex" alignItems="flex-start" gap={1}>
                    <Typography variant="h6">{title}</Typography>
                    <IconButton onClick={() => onEditToggle(true)} size="small">
                        <EditIcon fontSize="small" />
                    </IconButton>
                    {/* Кнопка ещё */}
                    <IconButton onClick={onMoreClick} size="small">
                        <MoreVertIcon fontSize="small" />
                    </IconButton>
                </Box>
            ) : (
                <TextField
                    label="Заголовок"
                    value={title}
                    inputRef={titleRef}
                    onChange={(e) => onChange(e.target.value)}
                    fullWidth
                />
            )}
        </Box>
    )
}
