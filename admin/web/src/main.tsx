import React from 'react'; import {createRoot} from 'react-dom/client'; import {RouterProvider} from 'react-router-dom'; import {router} from './app/router'; import './styles.css';
createRoot(document.getElementById('root')!).render(<React.StrictMode><RouterProvider router={router}/></React.StrictMode>);
