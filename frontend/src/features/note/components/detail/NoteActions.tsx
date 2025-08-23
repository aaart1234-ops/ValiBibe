import { Box, Button, CircularProgress } from '@mui/material'
import DeleteIcon from '@mui/icons-material/Delete'
import ArchiveIcon from '@mui/icons-material/Archive'

interface NoteActionsProps {
    isEditing: boolean
    noteArchived: boolean
    isSaving: boolean
    isUnarchiving: boolean
    onSave: () => void
    onCancelEdit: () => void
    onDelete: () => void
    onArchive: () => void
    onUnarchive: () => void
}

export const NoteActions = ({
                                isEditing,
                                noteArchived,
                                isSaving,
                                isUnarchiving,
                                onSave,
                                onCancelEdit,
                                onDelete,
                                onArchive,
                                onUnarchive,
                            }: NoteActionsProps) => {
    return isEditing ? (
        <Box display="flex" gap={1}>
            <Button variant="contained" onClick={onSave} disabled={isSaving}>
                {isSaving ? <CircularProgress size={20} /> : 'Сохранить'}
            </Button>
            <Button variant="outlined" onClick={onCancelEdit}>
                Отмена
            </Button>
        </Box>
    ) : (
        <Box display="flex" gap={1}>
            <Button variant="outlined" color="error" startIcon={<DeleteIcon />} onClick={onDelete}>
                Удалить
            </Button>
            {!noteArchived && (
                <Button variant="outlined" startIcon={<ArchiveIcon />} onClick={onArchive}>
                    Архивировать
                </Button>
            )}
            {noteArchived && (
                <Button variant="contained" color="secondary" onClick={onUnarchive} disabled={isUnarchiving}>
                    {isUnarchiving ? 'Восстановление...' : 'Из архива'}
                </Button>
            )}
        </Box>
    )
}
