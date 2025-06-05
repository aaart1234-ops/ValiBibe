import React from 'react'
import ReactDOM from 'react-dom/client'
import { Provider } from 'react-redux'
import { store } from './app/store'
import App from './App' // ✅ импортируем App.tsx

ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
        <Provider store={store}>
            <App /> {/* теперь всё идёт отсюда */}
        </Provider>
    </React.StrictMode>
)
