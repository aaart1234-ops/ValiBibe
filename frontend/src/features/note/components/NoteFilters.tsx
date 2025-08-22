import React from 'react'
import {
    Box, IconButton, Typography, Select, MenuItem, FormControl,
    InputLabel, TextField, Button, Tooltip, Fab, useMediaQuery
} from '@mui/material'
import { useTheme } from '@mui/material/styles'
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward'
import ArrowUpwardIcon from '@mui/icons-material/ArrowUpward'
import ViewModuleIcon from '@mui/icons-material/ViewModule'
import ViewListIcon from '@mui/icons-material/ViewList'
import AddIcon from '@mui/icons-material/Add'
import ArchiveIcon from '@mui/icons-material/Inventory2'
import ArchiveOutlinedIcon from '@mui/icons-material/Inventory2Outlined'

type SortBy = 'created_at' | 'next_review_at'
type SortDir = 'asc' | 'desc'
type ViewMode = 'card' | 'list'

type Props = {
    title?: string
    searchQuery: string
    onSearchChange: (v: string) => void
    sortBy: SortBy
    onSortByChange: (v: SortBy) => void
    sortDirection: SortDir
    onToggleSortDirection: () => void
    viewMode: ViewMode
    onToggleViewMode: () => void
    onOpenCreateDialog: () => void
}

const NoteFilters: React.FC<Props> = ({
                                          searchQuery,
                                          onSearchChange,
                                          sortBy,
                                          onSortByChange,
                                          sortDirection,
                                          onToggleSortDirection,
                                          viewMode,
                                          onToggleViewMode,
                                          onOpenCreateDialog,
                                      }) => {
    const theme = useTheme()
    const isMobile = useMediaQuery(theme.breakpoints.down('sm'))

    return (
        <Box mt={4} mb={2} sx={{ pl: isMobile ? 1 : 4, pr: isMobile ? 1 : 4 }}>
            <Box display="flex" flexDirection="column" gap={2}>
                <Box
                    display="flex"
                    flexDirection={isMobile ? 'column' : 'row'}
                    justifyContent="space-between"
                    gap={0}
                    flexWrap="wrap"
                >
                    <Box
                        display="flex"
                        flexDirection="row"
                        justifyContent="space-between"
                        gap={2}
                        flexWrap="wrap"
                        flexGrow={1}
                        minWidth={isMobile ? '100%' : 'auto'}
                    >
                        <TextField
                            label="Поиск по заметкам"
                            variant="outlined"
                            size="small"
                            value={searchQuery}
                            onChange={(e) => onSearchChange(e.target.value)}
                            sx={{ width: isMobile ? '100%' : 280 }}
                        />

                        <FormControl size="small" sx={{ minWidth: 160 }}>
                            <InputLabel id="sort-select-label">Сортировка</InputLabel>
                            <Select
                                labelId="sort-select-label"
                                value={sortBy}
                                label="Сортировка"
                                variant="outlined"
                                onChange={(e) => onSortByChange(e.target.value as SortBy)}
                            >
                                <MenuItem value="created_at">По дате создания</MenuItem>
                                <MenuItem value="next_review_at">По дате следующего повторения</MenuItem>
                            </Select>
                        </FormControl>

                        <Tooltip title={`Сортировать по ${sortDirection === 'asc' ? 'возрастанию' : 'убыванию'}`}>
                            <IconButton onClick={onToggleSortDirection}>
                                {sortDirection === 'asc' ? <ArrowUpwardIcon /> : <ArrowDownwardIcon />}
                            </IconButton>
                        </Tooltip>

                        <Tooltip title={viewMode === 'card' ? 'Список' : 'Карточки'}>
                            <IconButton onClick={onToggleViewMode}>
                                {viewMode === 'card' ? <ViewListIcon /> : <ViewModuleIcon />}
                            </IconButton>
                        </Tooltip>
                    </Box>

                    {!isMobile && (
                        <IconButton onClick={onOpenCreateDialog} color="primary" size="large" sx={{ padding: 0 }}>
                            <AddIcon fontSize="large" />
                        </IconButton>
                    )}
                </Box>

                {isMobile && (
                    <Fab color="primary" onClick={onOpenCreateDialog} sx={{ position: 'fixed', bottom: 24, right: 24, zIndex: 10 }}>
                        <AddIcon />
                    </Fab>
                )}
            </Box>
        </Box>
    )
}

export default NoteFilters
