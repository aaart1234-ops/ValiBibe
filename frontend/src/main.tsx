import React from 'react'
import ReactDOM from 'react-dom/client'
import { Provider } from 'react-redux'
import { store } from './app/store'
import App from './App' // ✅ импортируем App.tsx
import { MantineProvider } from '@mantine/core'
import '@mantine/core/styles.css' // ✅ ОБЯЗАТЕЛЬНО
import '@mantine/tiptap/styles.css' // ✅ ДЛЯ редактора

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <Provider store={store}>
            <MantineProvider>
                <App /> {/* теперь всё идёт отсюда */}
            </MantineProvider>
        </Provider>
    </React.StrictMode>
)
