import {
    Menu,
    MenuItem,
    ListItemIcon,
    ListItemText,
} from '@mui/material'
import ArchiveIcon from '@mui/icons-material/Archive'
import UnarchiveIcon from '@mui/icons-material/Unarchive'
import DeleteIcon from '@mui/icons-material/Delete'
import ContentCopyIcon from '@mui/icons-material/ContentCopy'
import FolderIcon from '@mui/icons-material/Folder'

interface NoteMoreMenuProps {
    anchorEl: HTMLElement | null
    open: boolean
    onClose: () => void
    noteArchived: boolean
    onArchive: () => void
    onUnarchive: () => void
    onDelete: () => void
    onDuplicate: () => void
    onMove: () => void
}

export const NoteMoreMenu = ({
                                 anchorEl,
                                 open,
                                 onClose,
                                 noteArchived,
                                 onArchive,
                                 onUnarchive,
                                 onDelete,
                                 onDuplicate,
                                 onMove,
                             }: NoteMoreMenuProps) => {
    return (
        <Menu
            anchorEl={anchorEl}
            open={open}
            onClose={onClose}
            anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
            transformOrigin={{ vertical: 'top', horizontal: 'right' }}
        >
            {!noteArchived ? (
                <MenuItem
                    onClick={() => {
                        onClose()
                        onArchive()
                    }}
                >
                    <ListItemIcon>
                        <ArchiveIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText primary="Архивировать" />
                </MenuItem>
            ) : (
                <MenuItem
                    onClick={() => {
                        onClose()
                        onUnarchive()
                    }}
                >
                    <ListItemIcon>
                        <UnarchiveIcon fontSize="small" />
                    </ListItemIcon>
                    <ListItemText primary="Восстановить" />
                </MenuItem>
            )}

            <MenuItem
                onClick={() => {
                    onClose()
                    onDelete()
                }}
            >
                <ListItemIcon>
                    <DeleteIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Удалить" />
            </MenuItem>

            <MenuItem
                onClick={() => {
                    onClose()
                    onDuplicate()
                }}
            >
                <ListItemIcon>
                    <ContentCopyIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Дублировать" />
            </MenuItem>

            <MenuItem
                onClick={() => {
                    onClose()
                    onMove()
                }}
            >
                <ListItemIcon>
                    <FolderIcon fontSize="small" />
                </ListItemIcon>
                <ListItemText primary="Переместить" />
            </MenuItem>
        </Menu>
    )
}
