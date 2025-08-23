import { Box, Typography, TextField, IconButton } from '@mui/material'
import EditIcon from '@mui/icons-material/Edit'
import { RefObject } from 'react'

interface NoteTitleProps {
    title: string
    isEditing: boolean
    onChange: (value: string) => void
    onEditToggle: (val: boolean) => void
    titleRef: RefObject<HTMLInputElement>
}

export const NoteTitle = ({ title, isEditing, onChange, onEditToggle, titleRef }: NoteTitleProps) => {
    return (
        <Box sx={{ pl: 6 }}>
            {!isEditing ? (
                <Box display="flex" alignItems="center" gap={1}>
                    <Typography variant="h5">{title}</Typography>
                    <IconButton onClick={() => onEditToggle(true)} size="small">
                        <EditIcon fontSize="small" />
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
