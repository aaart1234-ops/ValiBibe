import React from 'react'
import { Box, FormControl, MenuItem, Pagination, Select } from '@mui/material'

type Props = {
    total: number
    limit: number
    page: number // нумерация с 0
    onLimitChange: (newLimit: number) => void
    onPageChange: (newPageZeroBased: number) => void
}

const NotePagination: React.FC<Props> = ({
                                             total,
                                             limit,
                                             page,
                                             onLimitChange,
                                             onPageChange,
                                         }) => {
    const pagesCount = Math.ceil(total / limit)

    return (
        <Box display="flex" justifyContent="center" alignItems="center" mt={8} gap={2}>
            <FormControl size="small">
                <Select
                    value={limit}
                    onChange={(e) => onLimitChange(Number(e.target.value))}
                >
                    <MenuItem value={5}>5</MenuItem>
                    <MenuItem value={10}>10</MenuItem>
                    <MenuItem value={20}>20</MenuItem>
                    <MenuItem value={50}>50</MenuItem>
                </Select>
            </FormControl>

            {total > limit && (
                <Pagination
                    count={pagesCount}
                    page={page + 1}
                    onChange={(_, value) => onPageChange(value - 1)}
                    color="primary"
                />
            )}
        </Box>
    )
}

export default NotePagination
